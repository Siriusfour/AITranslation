package RabbitMQ

import (
	"AITranslatio/Global"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	handller    func(message []byte)
	consumerNum int
	conn        *amqp.Connection
	durable     bool
	autoDelete  bool
	queueName   string
	status      byte
}

func (c *consumer) Received(callbackFunDealSmg func(receivedData string)) {
	defer c.close()

	for i := 1; i <= c.consumerNum; i++ {

		go func() {
			megs := c.createConsumer()
			meg := <-megs

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
