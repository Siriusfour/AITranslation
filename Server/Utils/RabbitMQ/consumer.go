package RabbitMQ

import (
	"AITranslatio/Global"
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Handler func(ctx context.Context, body []byte)

type Config struct {
	URI           string
	Queue         string
	Durable       bool //是否持久化
	prefetch      int  //每个消费者的最大未确认信息数量（背压）
	workers       int  //每个消费者的并发协程数量
	RetryBase     time.Duration
	RetryMaxCount time.Duration
	RetryWaitTime time.Duration
	EnableConfirm bool //是否开启发布确认
}

type Client struct {
	config Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	closed atomic.Bool

	mu           sync.Mutex
	notifyCloseC chan *amqp.Error
	wg           sync.WaitGroup
}

// Connect 创建连接，失败时重试，直到成功/超出最大时间或者次数
func (c *Client) Connect(ctx context.Context) error {

	//边界检查
	if c.config.RetryBase <= 0 {
		c.config.RetryBase = time.Second
	}
	if c.config.RetryWaitTime <= 60 {
		c.config.RetryWaitTime = 30 * time.Second
	}

	backOff := c.config.RetryBase
	tries := 0

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if c.closed.Load() {
			return errors.New("client closed")
		}

		//建立连接
		conn, err := amqp.Dial(c.config.URI)
		if err == nil {
			c.conn = conn
			ch, err := conn.Channel()
			if err != nil {
				return err
			}

			//开启发布确认
			if c.config.EnableConfirm {
				if err := ch.Confirm(false); err != nil {
					_ = ch.Close()
					_ = conn.Close()
					return err
				}
			} else {
				c.mu.Lock()
				c.conn = conn
				c.ch = ch
				c.notifyCloseC = make(chan *amqp.Error, 1)
				conn.NotifyClose(c.notifyCloseC)
				c.mu.Unlock()

				log.Println("[MQ] ✅ connected")
				// 启动连接监视器（断线重连）
				c.wg.Add(1)
				go c.watchReconnect()
				return nil
			}

		} else { //连接失败，开始重试

		}
	}

}

type consumer struct {
	handller    func(message []byte)
	consumerNum int
	conn        *amqp.Connection
	durable     bool
	autoDelete  bool
	queueName   string
	status      byte
}

func (c *consumer) Received(callbackFunDealSmg func(receivedData amqp.Delivery)) {
	defer c.close()
	for i := 1; i <= c.consumerNum; i++ {
		go func() {
			megs := c.createConsumer()
			meg := <-megs
			callbackFunDealSmg(meg)
		}()

	}
}

func (c consumer) createConsumer() <-chan amqp.Delivery {

	channel, err := c.conn.Channel()
	if err != nil {
		c.log(err)
		return nil
	}

	queue, err := channel.QueueDeclare(
		c.queueName,
		c.durable,
		true,
		false,
		false,
		nil,
	)

	msgs, err := channel.Consume(
		queue.Name,
		"",    //  消费者标记，请确保在一个消息通道唯一
		true,  //是否自动确认，这里设置为 true，自动确认
		false, //是否私有队列，false标识允许多个 consumer 向该队列投递消息，true 表示独占
		false, //RabbitMQ不支持noLocal标志。
		false, // 队列如果已经在服务器声明，设置为 true ，否则设置为 false；
		nil,
	)
	if err != nil {
		c.log(err)
		return nil
	}
	return msgs
}

func (c *consumer) close() {
	_ = c.conn.Close()
}

func (c *consumer) log(err error) {
	Global.Logger.Error(err.Error())
}
