package RabbitMQ

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"

	"time"
)

// Publish  投放消息
func (c *Client) Publish(exchange, routingKey string, body []byte) error {

	c.mu.RLock()
	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("conncet创建channel失败：%w", err)
	}
	defer ch.Close()
	c.mu.RUnlock()

	// 🔴监听 Channel 关闭原因
	closeChan := make(chan *amqp.Error, 1)
	ch.NotifyClose(closeChan)

	// 开启 Confirm
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("enable confirm failed: %w", err)
	}
	acks := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// Publish ... (你的原有代码)
	if err := ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/octet-stream",
		Body:         body,
		MessageId:    "DEBUG_ID", // 暂时随便写
	}); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	// 等待结果
	select {
	case conf, ok := <-acks:
		// 如果 ok 为 false，说明 channel 被关闭了
		if !ok {
			// 尝试读取关闭原因
			select {
			case reason := <-closeChan:
				return fmt.Errorf("Channel被关闭，原因: %v", reason)
			default:
				return errors.New("Channel被异常关闭，且无具体原因")
			}
		}
		if !conf.Ack {
			// 这里才是真正的 NACK (资源不足等)
			// 顺便看看有没有关闭错误
			select {
			case reason := <-closeChan:
				return fmt.Errorf("收到NACK，且Channel关闭: %v", reason)
			default:
				return errors.New("收到NACK (可能是磁盘满或队列溢出)")
			}
		}
	case <-time.After(5 * time.Second):
		return errors.New("publish超时")
	}

	return nil
}
