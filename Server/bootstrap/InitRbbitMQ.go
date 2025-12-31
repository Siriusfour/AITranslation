package bootstrap

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Utils/RabbitMQ"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

func DialMQ(URI string, cfg amqp.Config) (*amqp.Connection, error) {
	return amqp.DialConfig(URI, cfg)
}

func InitMQClient(cfg interf.ConfigInterface, logger *zap.Logger) *RabbitMQ.Client {

	config := &RabbitMQ.Config{
		URI:     cfg.GetString("RabbitMq.WorkQueue.Addr"),
		Durable: cfg.GetBool("RabbitMq.WorkQueue.Durable"),
		Workers: cfg.GetInt("RabbitMq.WorkQueue.Workers"),

		RetryMaxTries: cfg.GetInt("RabbitMq.WorkQueue.Workers"),
		RetryBaseTime: time.Duration(cfg.GetInt("RabbitMq.WorkQueue.RetryBaseTime")), //ms
		RetryWaitTime: time.Duration(cfg.GetInt("RabbitMq.WorkQueue.Workers")),
		EnableConfirm: cfg.GetBool("RabbitMq.WorkQueue.EnableConfirm"), //是否开启发布确认

		// 死信队列
		DLXExchange:   cfg.GetString("RabbitMq.WorkQueue.DLXExchange"),   // e.g. "dlx.exchange"
		DLXRoutingKey: cfg.GetString("RabbitMq.WorkQueue.DLXRoutingKey"), // e.g. "dlx.key"
		DLQName:       cfg.GetString("RabbitMq.WorkQueue.DLQName"),       // e.g. "dlq"

		//生产确认，消费确认
		ConfirmTimeout: time.Duration(cfg.GetInt64("RabbitMq.WorkQueue.ConfirmTimeout")), // e.g. 5s
		MaxPubRetries:  cfg.GetInt("RabbitMq.WorkQueue.MaxPubRetries"),                   // e.g. 3

		// 延迟重发队列配置项
		RetryExchange: cfg.GetString("RabbitMq.WorkQueue.RetryExchange"), // e.g. "retry.exchange"
		RetryTime:     time.Duration(cfg.GetInt("RabbitMq.WorkQueue.RetryTime")),
	}

	client := &RabbitMQ.Client{
		Wg:     sync.WaitGroup{},
		Cfg:    cfg,
		Logger: logger,
		Config: *config,
		Conn:   nil,
		Closed: atomic.Bool{},
		DialMQ: DialMQ,
	}

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	return client

}
