package RabbitMQ

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// Publish  投放消息

func (c *Client) Publish(queueName string, body []byte, exchangeType string) error {

	c.mu.RLock()
	ch := c.ch
	c.mu.RUnlock()
	if ch == nil {
		return errors.New("channel not available")
	}

	args := amqp.Table{
		// 消息10秒没被消费就过期，成为“死信”
		"x-message-ttl": int32(10000), // 10 秒
		// 过期的消息会转发到以下 DLX
		"x-dead-letter-exchange":    "dlx.exchange",
		"x-dead-letter-routing-key": "dlx.key",
	}

	if _, err := ch.QueueDeclare(queueName, false, false, false, false, args); err != nil {
		return fmt.Errorf("MQ创建/链接队列失败: %w", err)
	}
	// 选择器：如果启用发布确认，注册监听确认通道
	var acks <-chan amqp.Confirmation
	if c.config.EnableConfirm {
		acks = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	if err := ch.Publish("", queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/octet-stream",
		Body:         body,
	}); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}
	if c.config.EnableConfirm {
		select {
		case conf := <-acks:
			if !conf.Ack {
				return errors.New("publish not acknowledged by broker")
			}
		case <-time.After(5 * time.Second):
			return errors.New("publish confirm timeout")
		}
	}
	return nil
}
