package RabbitMQ

import (
	"AITranslatio/Global"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type Handler func(body []byte) error

type Config struct {
	URI            string
	Durable        bool //是否持久化
	prefetch       int  //每个消费者的最大未确认信息数量（背压）
	Workers        int  //每个消费者的并发协程数量
	Heartbeat      time.Duration
	ConnectionName string

	//掉线重连配置项
	RetryMaxTries int
	RetryBaseTime time.Duration
	RetryWaitTime time.Duration
	EnableConfirm bool //是否开启发布确认

	// 死信队列
	DLXExchange   string // e.g. "dlx.exchange"
	DLXRoutingKey string // e.g. "dlx.key"
	DLQName       string // e.g. "dlq"

	//生产确认，消费确认
	ConfirmTimeout time.Duration // e.g. 5s
	MaxPubRetries  int           // e.g. 3

	// 演出重发队列配置项
	RetryExchange string // e.g. "重发交换机"
	RetryTime     time.Duration
}

type Client struct {
	wg            sync.WaitGroup
	Config        Config
	mu            sync.RWMutex
	Conn          *amqp.Connection
	Closed        atomic.Bool
	notifychannel chan *amqp.Error
}

func InitClient() (client *Client, err error) {

	config := &Config{
		URI:     Global.Config.GetString("RabbitMq.WorkQueue.Addr"),
		Durable: Global.Config.GetBool("RabbitMq.WorkQueue.Durable"),
		Workers: Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),

		RetryMaxTries: Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),
		RetryBaseTime: time.Duration(Global.Config.GetInt("RabbitMq.WorkQueue.RetryBaseTime")), //ms
		RetryWaitTime: time.Duration(Global.Config.GetInt("RabbitMq.WorkQueue.Workers")),
		EnableConfirm: Global.Config.GetBool("RabbitMq.WorkQueue.EnableConfirm"), //是否开启发布确认

		// 死信队列
		DLXExchange:   Global.Config.GetString("RabbitMq.WorkQueue.DLXExchange"),   // e.g. "dlx.exchange"
		DLXRoutingKey: Global.Config.GetString("RabbitMq.WorkQueue.DLXRoutingKey"), // e.g. "dlx.key"
		DLQName:       Global.Config.GetString("RabbitMq.WorkQueue.DLQName"),       // e.g. "dlq"

		//生产确认，消费确认
		ConfirmTimeout: time.Duration(Global.Config.GetInt64("RabbitMq.WorkQueue.ConfirmTimeout")), // e.g. 5s
		MaxPubRetries:  Global.Config.GetInt("RabbitMq.WorkQueue.MaxPubRetries"),                   // e.g. 3

		// 演出重发队列配置项
		RetryExchange: Global.Config.GetString("RabbitMq.WorkQueue.RetryExchange"), // e.g. "retry.exchange"
		RetryTime:     time.Duration(Global.Config.GetInt("RabbitMq.WorkQueue.RetryTime")),
	}

	client = &Client{
		wg:     sync.WaitGroup{},
		Config: *config,
		Conn:   nil,
		Closed: atomic.Bool{},
	}

	//安全检查，判断队列是否齐全

	err = client.Connect()
	err = client.EnsureDLX()
	err = client.EnsureQueue()

	if err != nil {
		return nil, fmt.Errorf("mq客户端初始化失败: %w", err)
	}

	return client, nil

}

// client的客户端连接到MQ服务器
func (c *Client) Connect() error {
	// 心跳默认值
	hb := c.Config.Heartbeat
	if hb <= 0 {
		hb = 10 * time.Second
	}

	// 退避基础时间 & 上限
	base := c.Config.RetryBaseTime
	if base <= 0 {
		base = time.Second
	}
	maxWait := c.Config.RetryWaitTime
	if maxWait <= 0 {
		maxWait = 30 * time.Second
	}

	cfg := amqp.Config{
		Heartbeat: hb,
		Properties: amqp.Table{
			"connection_name": c.Config.ConnectionName,
		},
	}

	backOff := base
	tries := 0

	for {
		if c.Closed.Load() {
			return errors.New("client closed")
		}

		// 2. 尝试建立连接
		conn, err := amqp.DialConfig(c.Config.URI, cfg)
		if err == nil { //连接成功

			c.mu.Lock()

			old := c.Conn
			if old != nil {
				_ = old.Close()
			} //关闭旧链接

			c.notifychannel = make(chan *amqp.Error, 1) //,注册notifyCh
			c.Conn.NotifyClose(c.notifychannel)
			c.Conn = conn
			go c.watchReconnect()
			c.mu.Unlock()

			Global.Logger["MQ"].Info("[MQ] ✅ connected", zap.String("uri", c.Config.URI))
			return nil
		}

		// 3. 连接失败，判断重试次数
		tries++
		if c.Config.RetryMaxTries > 0 && tries >= c.Config.RetryMaxTries {
			Global.Logger["MQ"].Error("MQ 建立连接失败，超过最大重试次数",
				zap.Error(err),
				zap.Int("tries", tries),
			)
			return err
		}

		// 4. 退避等待
		timer := time.NewTimer(backOff)
		select {
		case <-timer.C:
		}

		// 指数退避并封顶
		backOff *= 2
		if backOff > maxWait {
			backOff = maxWait
		}
	}
}

// watchReconnect 只负责：
// 1. 监听当前连接的关闭事件
// 2. 触发重连（带重试、sleep）
// 3. 重连成功后重新注册 NotifyClose，然后再起一个新的 watcher
func (c *Client) watchReconnect() {
	for {
		select {
		case err := <-c.notifychannel:
			Global.Logger["MQ"].Error("MQ 连接关闭", zap.Error(err))

			// 清空当前连接（加锁保护）
			c.mu.Lock()
			old := c.Conn
			c.Conn = nil
			c.mu.Unlock()
			if old != nil {
				_ = old.Close()
			}

			// 开始重连，带简单的重试+sleep
			for {

				if c.Closed.Load() {
					return
				}

				connErr := c.Connect()
				if connErr == nil {
					// 重连成功：给新连接注册 NotifyClose，并重新起一个 watcher
					c.mu.RLock()
					newConn := c.Conn
					c.mu.RUnlock()

					notify := make(chan *amqp.Error, 1)
					newConn.NotifyClose(notify)

					c.mu.Lock()
					c.notifychannel = notify
					c.mu.Unlock()

					// 起一个新的 watcher 监听新连接
					go c.watchReconnect()
					return
				}

				Global.Logger["MQ"].Error("MQ 重连失败，将在 5s 后重试", zap.Error(connErr))
				time.Sleep(5 * time.Second)
			}

		}
	}
}

// 声明死信队列和死信交换机
func (c *Client) EnsureDLX() error {
	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 声明死信交换机（DLX）
	err = ch.ExchangeDeclare(
		c.Config.DLXExchange, // 交换机名
		"direct",             // 类型：一般 direct 即可
		true,                 // durable 持久化
		false,                // autoDelete
		false,                // internal
		false,                // noWait
		nil,                  // args
	)
	if err != nil {
		return err
	}

	dlq, err := ch.QueueDeclare(
		"dlq", // 死信队列名称
		true,  // durable 持久化
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}
	_ = dlq // 视情况使用

	err = ch.QueueBind(
		"dlq",                  // 队列名
		c.Config.DLXRoutingKey, // routing key
		"dlx.exchange",         // 交换机名
		false,                  // noWait
		nil,                    // args
	)
	if err != nil {
		return err
	}

	return nil
}

// 由配置文件声明业务队列，并挂上 DLX 参数
func (c *Client) EnsureQueue() error {

	exchanges := Global.Config.Get("RabbitMq.WorkQueue.exchanges").([]interface{})
	fmt.Println(exchanges)
	for _, queueName := range exchanges {
		fmt.Println(queueName)
		ch, err := c.Conn.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		args := amqp.Table{}
		if c.Config.DLXExchange != "" {
			args["x-dead-letter-exchange"] = c.Config.DLXExchange
		}
		if c.Config.DLXRoutingKey != "" {
			args["x-dead-letter-routing-key"] = c.Config.DLXRoutingKey
		}

		_, err = ch.QueueDeclare(
			"",
			c.Config.Durable,
			false, // autoDelete
			false, // exclusive
			false, // noWait
			args,
		)
	}

	return nil
}

// EnsureRetryTopology 为某个业务队列声明“只有一个 5 分钟延迟的重试队列”
// 延迟队列名：retry.<queue>.5m
// 逻辑：消费失败 -> 发送到该延迟队列；消息在队列里等待 5 分钟，TTL 到期后经 DLX 回到原业务队列。
func (c *Client) EnsureRetryTopology(queue string) error {

	ch, err := c.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 1. 声明重试交换机
	if err := ch.ExchangeDeclare(
		c.Config.RetryExchange,
		"direct",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	); err != nil {
		return err
	}

	// 2. 声明一个 2 分钟延迟的重试队列
	delay := 2 * time.Minute
	retryQueueName := "retry." + queue + ".2m"

	args := amqp.Table{
		// 消息在重试队列中的存活时间：2 分钟
		"x-message-ttl": int(delay / time.Millisecond),
		// TTL 到期后，通过默认交换机回到原业务队列，所有队列都自动绑定到默认交换机 ， routing key = 队列名
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": queue,
	}

	if _, err := ch.QueueDeclare(
		retryQueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		args,
	); err != nil {
		return err
	}

	// 3. 绑定：retry.<queue>.2m
	if err := ch.QueueBind(
		retryQueueName,
		retryQueueName, // routingKey
		c.Config.RetryExchange,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}
