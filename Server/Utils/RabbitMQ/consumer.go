package RabbitMQ

import (
	"AITranslatio/Global"
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Handler func(ctx context.Context, body []byte) error

type Config struct {
	URI           string
	Durable       bool //是否持久化
	prefetch      int  //每个消费者的最大未确认信息数量（背压）
	workers       int  //每个消费者的并发协程数量
	RetryBase     time.Duration
	RetryMaxTries int
	RetryWaitTime time.Duration
	EnableConfirm bool //是否开启发布确认
}

type Client struct {
	config Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	closed atomic.Bool

	mu           sync.RWMutex
	notifyCloseC chan *amqp.Error
	wg           sync.WaitGroup
	dealQueue    chan amqp.Delivery
}

func InitClient() (*Client, error) {

	config := &Config{
		URI:       Global.Config.GetString("RabbitMq.WorkQueue.Addr"),
		Durable:   Global.Config.GetBool("RabbitMq.WorkQueue.Durable"),
		workers:   Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),
		RetryBase: Global.Config.GetDuration("RabbitMq.WorkQueue.RetryBase"),
	}

	conn, err := amqp.Dial()
	if err != nil {
		Global.Logger["MQ"].Error("Failed to connect RabbitMQ" + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		Global.Logger["MQ"].Error("Failed to connect RabbitMQ" + err.Error())
	}

	// 1️⃣ 声明死信交换机
	err = ch.ExchangeDeclare(
		"dlx.exchange", // 死信交换机名
		"direct",       // 类型 direct
		true,           // durable
		false,          // autoDelete
		false,          // internal
		false,          // noWait
		nil,            // args
	)

	//声明死信队列
	_, err = ch.QueueDeclare(
		"dead_letter_queue", // 死信队列名称
		true,                // durable
		false,               // autoDelete
		false,               // exclusive
		false,               // noWait
		nil,                 // args
	)

	//二者绑定
	err = ch.QueueBind(
		"dead_letter_queue", // 队列名
		"dlx.key",           // 路由键
		"dlx.exchange",      // 绑定的死信交换机
		false,
		nil,
	)

}

// Connect 创建连接，失败时重试，直到成功/超出最大时间和次数
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

			//发布确认处理
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
			tries++
			if c.config.RetryMaxTries > 0 && tries > c.config.RetryMaxTries {
				Global.Logger["MQ"].Error("connect retries exceeded: %w", zap.Error(err))
				return err
			}

			timer := time.NewTimer(backOff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
			}
			backOff *= 2
			if backOff > c.config.RetryWaitTime {
				backOff = c.config.RetryWaitTime
			}

			Global.Logger["MQ"].Error("[MQ] ❌ connect is failed :%w", zap.Error(err))
		}
	}
}

func (c *Client) watchReconnect() {

	defer c.wg.Done()

	c.mu.Lock()
	errCh := c.notifyCloseC
	c.mu.Unlock()

	for {
		if errCh == nil {
			return
		}

		select {
		case <-errCh:
			//释放旧资源 ch，conn
			c.mu.Lock()
			if c.ch != nil {
				_ = c.ch.Close()
				c.ch = nil
			}
			if c.conn != nil {
				_ = c.conn.Close()
				c.conn = nil
			}
			c.notifyCloseC = nil
			// 重连（无限重试，直到成功或外部关闭）
			c.mu.Unlock()
			ctx := context.Background()
			if err := c.Connect(ctx); err == nil {
				Global.Logger["MQ"].Info("MQ is re connected", zap.String("MQ", c.config.URI))
				return
			}

		case <-time.After(50 * time.Second):
		}
	}

}

func (c *Client) Consume(ctx context.Context, queueName string, handler Handler) error {

	if queueName == "" {
		return errors.New("queue name required")
	}

	workers := c.config.workers
	if workers <= 0 {
		workers = 1
	}

	for i := 0; i < workers; i++ {
		c.wg.Add(1)
		go func(idx int) {
			defer c.wg.Done()
			for {
				// 优雅地处理取消信号
				select {
				case <-ctx.Done():
					log.Printf("[MQ] worker=%d context cancelled", idx)
					return
				default:
				}

				if c.closed.Load() {
					log.Printf("[MQ] worker=%d channel closed", idx)
					return
				}

				c.mu.RLock()
				ch := c.ch
				c.mu.RUnlock()
				if ch == nil {
					log.Printf("[MQ] worker=%d channel is nil, retrying", idx)
					time.Sleep(300 * time.Millisecond)
					continue
				}

				prefetch := c.config.prefetch
				if prefetch <= 0 {
					prefetch = 50
				}
				_ = ch.Qos(prefetch, 0, false)

				if _, err := ch.QueueDeclare(queueName, c.config.Durable, false, false, false, nil); err != nil {
					log.Printf("[MQ] worker=%d declare queue error: %v", idx, err)
					time.Sleep(time.Second)
					continue
				}

				msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
				if err != nil {
					log.Printf("[MQ] worker=%d consume error: %v", idx, err)
					time.Sleep(time.Second)
					continue
				}

				for msg := range msgs {
					func() {
						defer func() {
							if r := recover(); r != nil {
								log.Printf("[MQ] worker=%d panic: %v", idx, r)
								_ = msg.Nack(false, true)
							}
						}()

						if err := handler(ctx, msg.Body); err != nil {
							log.Printf("[MQ] worker=%d handler error: %v", idx, err)
							_ = msg.Nack(false, true) // 可替换为死信队列处理
							return
						}

						_ = msg.Ack(false)
					}()
				}
			}
		}(i)
	}
	return nil
}
