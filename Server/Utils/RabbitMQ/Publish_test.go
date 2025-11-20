package RabbitMQ

import "testing"

//TODO 调用MQ生产，消费逻辑,看是否正确

func TestPublish(t *testing.T) {

	c := initTestClient(mockDialSuccess)

	err := c.Publish("user.websocket", "Message", []byte("Message"))
	if err != nil {
		return
	}

}
