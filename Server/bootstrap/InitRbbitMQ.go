package bootstrap

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/RabbitMQ"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"sync/atomic"
	"time"
)

func DialMQ(URI string, cfg amqp.Config) (*amqp.Connection, error) {
	return amqp.DialConfig(URI, cfg)
}

func InitMQClient() {

	config := &RabbitMQ.Config{
		URI:     Global.Config.GetString("RabbitMq.WorkQueue.Addr"),
		Durable: Global.Config.GetBool("RabbitMq.WorkQueue.Durable"),
		Workers: Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),

		RetryMaxTries: Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),
		RetryBaseTime: time.Duration(Global.Config.GetInt("RabbitMq.WorkQueue.RetryBaseTime")), //ms
		RetryWaitTime: time.Duration(Global.Config.GetInt("RabbitMq.WorkQueue.Workers")),
		EnableConfirm: Global.Config.GetBool("RabbitMq.WorkQueue.EnableConfirm"), //是否开启发布确认

		// 死信队列
		DLXExchange:   Global.Config.GetString("RabbitMq.WorkQueue.DLXExchange"),   // e.g. "dlx.exchange"
		DLXRoutingKey: Global.Config.GetString("RabbitMq.WorkQueue.DLXRoutingKey"), // e.g. "dlx.key"
		DLQName:       Global.Config.GetString("RabbitMq.WorkQueue.DLQName"),       // e.g. "dlq"

		//生产确认，消费确认
		ConfirmTimeout: time.Duration(Global.Config.GetInt64("RabbitMq.WorkQueue.ConfirmTimeout")), // e.g. 5s
		MaxPubRetries:  Global.Config.GetInt("RabbitMq.WorkQueue.MaxPubRetries"),                   // e.g. 3

		// 延迟重发队列配置项
		RetryExchange: Global.Config.GetString("RabbitMq.WorkQueue.RetryExchange"), // e.g. "retry.exchange"
		RetryTime:     time.Duration(Global.Config.GetInt("RabbitMq.WorkQueue.RetryTime")),
	}

	client := &RabbitMQ.Client{
		Wg:     sync.WaitGroup{},
		Config: *config,
		Conn:   nil,
		Closed: atomic.Bool{},
		DialMQ: DialMQ,
	}

	Global.RabbitmqClient = client

	//err := client.Connect()
	//if err != nil {
	//	panic(err)
	//}

}
