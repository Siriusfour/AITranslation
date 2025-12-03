package RabbitMQ

import (
	"AITranslatio/Utils/SnowFlak"
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

	// 如果启用发布确认，注册监听确认通道
	var acks <-chan amqp.Confirmation

	if c.Config.EnableConfirm {
		if err := ch.Confirm(false); err != nil {
			return fmt.Errorf("enable confirm failed: %w", err)
		}
		acks = ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	//Publish向队列推送消息
	if err := ch.Publish(exchange, routingKey, false, false, amqp.Publishing{

		DeliveryMode: amqp.Persistent,
		ContentType:  "application/octet-stream",
		Body:         body,
		MessageId:    SnowFlak.CreateSnowflakeFactory(c.Cfg, c.Logger).GetIDString(),
	}); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	//推送消息后监测5秒内有没有来自MQ服务器的ACK,没有的话返回错误
	if c.Config.EnableConfirm {
		select {
		case conf := <-acks:
			if !conf.Ack {
				return errors.New("推送信息时没有收到来自broke的ACK")
			}
		case <-time.After(5 * time.Second):
			return errors.New("publish超时")
		}
	}
	return nil
}
