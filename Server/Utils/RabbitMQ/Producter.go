package RabbitMQ

import (
	"AITranslatio/Global"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	connect    *amqp.Connection //	地址
	queueName  string           //队列名称
	durable    bool
	occurError error
}

// CreateProducer  创建一个生产者
func CreateProducer() (*Producer, error) {
	// 获取配置信息
	conn, err := amqp.Dial(Global.Config.GetString("RabbitMq.WorkQueue.Addr"))
	queueName := Global.Config.GetString("RabbitMq.WorkQueue.QueueName")
	durable := Global.Config.GetBool("RabbitMq.WorkQueue.Durable")

	if err != nil {
		Global.Logger.Error(err.Error())
		return nil, err
	}

	prod := &Producer{
		connect:   conn,
		queueName: queueName,
		durable:   durable,
	}
	return prod, nil
}
