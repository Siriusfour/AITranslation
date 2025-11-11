package RabbitMQ

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// 使用n个worker持续监听某个队列，并把消费失败的err写入ch中
func (c *Client) consumer(ctx context.Context, queueName string, handleFunc Handler, errChan chan error) {

	//监听调用者是否关闭consumer的持续监听

	for i := 1; i <= c.Config.Workers; i++ {

		go func(workerID int, ctx context.Context) {

			//创建ch,失败则重连
			for {
				//conn断开，尝试重连
				if c.Conn == nil || c.Conn.IsClosed() {
					time.Sleep(10 * time.Second)
					continue
				}

				ch, err := c.Conn.Channel()
				if err != nil {
					errChan <- fmt.Errorf("Channel创建失败: %w", err)
					return
				}
				notifyClose := ch.NotifyClose(make(chan *amqp.Error))
				ch.Qos(c.Config.prefetch, 0, false)

				msgs, err := ch.Consume(
					queueName, // 队列名称
					"",        //  消费者标记，请确保在一个消息频道唯一
					false,     //是否自动确认，这里设置为 true，自动确认
					false,     //是否私有队列，false标识允许多个 consumer 向该队列投递消息，true 表示独占
					false,     //RabbitMQ不支持noLocal标志。
					false,     // 队列如果已经在服务器声明，设置为 true ，否则设置为 false；
					nil,
				)

				if err == nil {

				consumeLoop:
					for {
						select {
						case msg := <-msgs:
							err := handleFunc(msg.Body)
							if err != nil {
								msg.Nack(false, false)
								errChan <- fmt.Errorf("WorkerID:%v，消费失败：%w", workerID, err)
							} else {
								msg.Ack(false)
							}
						case <-ctx.Done():
							ch.Close()
							return

						case err := <-notifyClose:
							if err != nil {
								errChan <- fmt.Errorf("WorkerID:%v Channel 异常关闭：%v", workerID, err)
							}
							ch.Close()
							time.Sleep(5 * time.Second)
							break consumeLoop
						}
					}
				} else {
					errChan <- fmt.Errorf("consume调用失败：%w", err)
				}

			}

		}(i, ctx)

	}

	return
}

// 判断是否是有经过延迟重连队列
func (c *Client) IsRetryMessage() {}

// 投递到延迟重连队列
func (c *Client) SendToRetry() {}
