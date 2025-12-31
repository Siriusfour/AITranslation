package RabbitMQ

import (
	"AITranslatio/Config/interf"
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type Handler func(ctx context.Context, body []byte) error
type DialMQ func(string, amqp.Config) (*amqp.Connection, error)

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
	Wg            sync.WaitGroup
	Logger        *zap.Logger
	Cfg           interf.ConfigInterface
	mu            sync.RWMutex
	Config        Config
	Conn          *amqp.Connection
	Closed        atomic.Bool
	notifychannel chan *amqp.Error
	//封装的连接MQ服务器的函数，单元测试、集成测试、生产环境初始化时各自填充
	DialMQ DialMQ
}

// client的客户端连接到MQ服务器
func (c *Client) Connect() error {

	// 心跳默认值
	hb := c.Config.Heartbeat
	if hb <= 0 {
		hb = 10 * time.Second
	}

	// 退避基础时间 & 上限
	RetryBaseTime := c.Config.RetryBaseTime
	if RetryBaseTime <= 0 {
		RetryBaseTime = 100 * time.Millisecond
	}
	maxWait := c.Config.RetryWaitTime
	if maxWait <= 0 {
		maxWait = 15 * time.Second
	}

	cfg := amqp.Config{
		Heartbeat: hb,
		Properties: amqp.Table{
			"connection_name": c.Config.ConnectionName,
		},
	}

	backOff := RetryBaseTime
	tries := 0

	for {

		if c.Closed.Load() {
			return errors.New("client closed")
		}

		// 2. 尝试建立连接
		conn, err := c.DialMQ(c.Config.URI, cfg)
		if err == nil { //连接成功

			c.mu.Lock()
			old := c.Conn
			if old != nil {
				_ = old.Close()
			} //关闭旧链接
			c.Conn = conn
			c.notifychannel = make(chan *amqp.Error, 1) //,注册notifyCh

			fmt.Println("c.conn:", c.Conn)

			c.Conn.NotifyClose(c.notifychannel)
			c.Conn = conn
			go c.watchReconnect()
			c.mu.Unlock()

			c.Logger.Info("[MQ] ✅ connected", zap.String("uri", c.Config.URI))
			return nil
		}

		// 3. 连接失败，判断重试次数
		tries++
		if c.Config.RetryMaxTries > 0 && tries >= c.Config.RetryMaxTries {
			c.Logger.Error("MQ 建立连接失败，超过最大重试次数",
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
			c.Logger.Error("MQ 连接关闭", zap.Error(err))

			// 清空当前连接（加锁保护）
			c.mu.Lock()

			c.Conn = nil
			if c.notifychannel != nil {
				close(c.notifychannel)
			}
			c.mu.Unlock()

			// 开始重连，带简单的重试+sleep
			for {
				if c.Closed.Load() {
					return
				}

				connErr := c.Connect()
				if connErr == nil {
					// 重连成功：给新连接注册 NotifyClose，并break,重新监听
					c.mu.RLock()
					newConn := c.Conn
					c.mu.RUnlock()

					notify := make(chan *amqp.Error, 1)
					newConn.NotifyClose(notify)

					c.mu.Lock()
					c.notifychannel = notify
					c.mu.Unlock()

					// 返回外层循环，重新监听
					break

				}

				c.Logger.Error("MQ 重连失败，将在 5s 后重试", zap.Error(connErr))
				time.Sleep(5 * time.Second)
			}

		}
	}
}

func (c *Client) Close() {
	c.Closed.Store(true)
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
