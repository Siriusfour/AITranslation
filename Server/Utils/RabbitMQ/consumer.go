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
						//持续监听信息
						case msg := <-msgs:
							err := handleFunc(msg.Body)
							if err != nil {
								c.MessageFailHandler(msg, queueName)
								continue
							} else {
								msg.Ack(false)
							}
						//调用者主动关闭
						case <-ctx.Done():
							ch.Close()
							return
						//通道被broker关闭
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
