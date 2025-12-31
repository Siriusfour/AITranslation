package RabbitMQ

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func (c *Client) Consumer(ctx context.Context, queueName string, handleFunc Handler, errChan chan error) {

	// 启动 N 个 Worker
	for i := 1; i <= c.Config.Workers; i++ {

		go func(workerID int) {

			// 外部循环：负责断线重连
			for {
				// 0. 检查上下文是否已取消
				select {
				case <-ctx.Done():
					return
				default:
				}

				// 1. 获取连接 (如果连接断了，这里应该阻塞等待或者报错)
				// 建议把 Conn 的获取封装成一个阻塞直到可用的方法，或者在这里简单的 Sleep
				c.mu.RLock()
				conn := c.Conn
				c.mu.RUnlock()

				if conn == nil || conn.IsClosed() {
					time.Sleep(2 * time.Second) // 等待主重连协程恢复连接
					continue
				}

				// 2. 创建 Channel
				ch, err := conn.Channel()
				if err != nil {
					errChan <- fmt.Errorf("[Worker-%d] Channel创建失败: %v", workerID, err)
					time.Sleep(5 * time.Second) // 避退
					continue
				}

				// 3. 配置
				ch.Qos(c.Config.prefetch, 0, false)
				notifyClose := ch.NotifyClose(make(chan *amqp.Error, 1)) // 缓冲设为1更安全

				msgs, err := ch.Consume(
					queueName,
					fmt.Sprintf("worker-%d", workerID), // 唯一的 Tag
					false,                              // autoAck: false
					false,
					false,
					false,
					nil,
				)

				if err != nil {
					errChan <- fmt.Errorf("[Worker-%d] Consume失败: %v", workerID, err)
					ch.Close()
					time.Sleep(5 * time.Second)
					continue
				}

				// 4. 消费循环
				c.Logger.Info(fmt.Sprintf("[Worker-%d] 开始监听队列: %s", workerID, queueName))

			consumeLoop:
				for {
					select {
					case msg, ok := <-msgs:
						// ⚠️【关键修复】检查通道是否关闭
						if !ok {
							c.Logger.Warn(fmt.Sprintf("[Worker-%d] msgs通道已关闭，准备重连", workerID))
							break consumeLoop
						}

						// 执行业务逻辑
						err := handleFunc(ctx, msg.Body)
						if err != nil {
							// 失败处理：建议 Nack(false, true) 让消息重回队列，或者发死信
							c.MessageFailHandler(msg, queueName)
							// 注意：MessageFailHandler 内部必须做 Ack/Nack，否则消息会死锁
						} else {
							// 成功确认
							msg.Ack(false)
						}

					case <-ctx.Done():
						ch.Close()
						return // 退出 Worker

					case err := <-notifyClose:
						c.Logger.Error(fmt.Sprintf("[Worker-%d] Channel异常关闭: %v", workerID, err))
						break consumeLoop
					}
				}
			}
		}(i)
	}
}

func (c *Client) MessageFailHandler(msg amqp.Delivery, queueName string) {

	if c.IsRetryMessage(&msg) {
		err := c.SendTo(&msg, queueName, "deal")
		if err != nil {
			msg.Nack(false, true)

		}

	} else {
		err := c.SendTo(&msg, queueName, "retry")
		if err != nil {
			return
		}
	}

}

// 判断是否是有经过延迟重连队列
func (c *Client) IsRetryMessage(msg *amqp.Delivery) bool {

	val, ok := msg.Headers["x-retry"]
	if ok {
		return val == 1
	} else {
		return false
	}
}

// SendTo 投递到死信or延迟重试队列
func (c *Client) SendTo(msg *amqp.Delivery, queueName string, Type string) error {

	if msg.Headers == nil {
		msg.Headers = amqp.Table{}
	}

	if Type == "retry" {
		msg.Headers["x-retry"] = 1
	}

	Exchange := Type + queueName

	// 打开一个新的 channel
	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("sendTO创建ch失败 %w", err)
	}
	defer ch.Close()

	// 直接发布消息到死信交换机
	return ch.Publish(
		Exchange,  // 已存在的 DLX
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   msg.ContentType,
			Body:          msg.Body,
			DeliveryMode:  amqp.Persistent, // 持久化消息
			CorrelationId: msg.CorrelationId,
			Headers:       msg.Headers,
		},
	)
}
