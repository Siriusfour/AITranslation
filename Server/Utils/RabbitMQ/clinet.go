package RabbitMQ

import (
	"AITranslatio/Global"
	"context"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync/atomic"
	"time"
)

type Handler func(ctx context.Context, body []byte) error

type Config struct {
	URI            string
	Durable        bool //是否持久化
	prefetch       int  //每个消费者的最大未确认信息数量（背压）
	workers        int  //每个消费者的并发协程数量
	Heartbeat      time.Duration
	ConnectionName string

	//掉线重连配置项
	RetryMaxTries int
	RetryBaseTime time.Duration
	RetryWaitTime time.Duration
	EnableConfirm bool //是否开启发布确认

	// 死信队列
	DLXExchange   string // e.g. "dlx.exchange"
	DLXRoutingKey string // e.g. "dlx.key"
	DLQName       string // e.g. "dlq"

	//生产确认，消费确认
	ConfirmTimeout time.Duration // e.g. 5s
	MaxPubRetries  int           // e.g. 3

	// 演出重发队列配置项
	RetryExchange string // e.g. "重发交换机"
	RetryBuckets  time.Duration
}

type Client struct {
	config Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	closed atomic.Bool
}

func InitClient() (*Client, error) {

	config := &Config{
		URI:     Global.Config.GetString("RabbitMq.WorkQueue.Addr"),
		Durable: Global.Config.GetBool("RabbitMq.WorkQueue.Durable"),
		workers: Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),

		RetryMaxTries: Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),
		RetryBaseTime: time.Duration(Global.Config.GetInt("RabbitMq.WorkQueue.RetryBaseTime")), //ms
		RetryWaitTime: time.Duration(Global.Config.GetInt64("RabbitMq.WorkQueue.Workers")),
		EnableConfirm: Global.Config.GetBool("RabbitMq.WorkQueue.Workers"), //是否开启发布确认

		// 死信队列
		DLXExchange:   Global.Config.GetString("RabbitMq.WorkQueue.Workers"), // e.g. "dlx.exchange"
		DLXRoutingKey: Global.Config.GetString("RabbitMq.WorkQueue.Workers"), // e.g. "dlx.key"
		DLQName:       Global.Config.GetString("RabbitMq.WorkQueue.Workers"), // e.g. "dlq"

		//生产确认，消费确认
		ConfirmTimeout: time.Duration(Global.Config.GetInt64("RabbitMq.WorkQueue.Workers")), // e.g. 5s
		MaxPubRetries:  Global.Config.GetInt("RabbitMq.WorkQueue.Workers"),                  // e.g. 3

		// 演出重发队列配置项
		RetryExchange: Global.Config.GetString("RabbitMq.WorkQueue.Workers"), // e.g. "retry.exchange"
		RetryBuckets:  time.Duration(Global.Config.GetInt64("RabbitMq.WorkQueue.Workers")),
	}

	conn, err := amqp.Dial(config.URI)
	if err != nil {
		Global.Logger["MQ"].Error("Failed to connect RabbitMQ" + err.Error())
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		Global.Logger["MQ"].Error("Failed to connect RabbitMQ" + err.Error())
		return nil, err
	}

	client := &Client{
		config: *config,
		conn:   conn,
		ch:     ch,
		closed: atomic.Bool{},
	}

	return client, nil
}

func (c *Client) connect(ctx context.Context) error {

	hb := c.config.Heartbeat
	connCount := 0

	cfg := amqp.Config{Heartbeat: hb, Properties: amqp.Table{"connection_name": c.config.ConnectionName}}

	if ctx.Err() != nil || c.closed.Load() {
		return errors.New("context deadline exceeded")
	}

	conn, err := amqp.DialConfig(c.config.URI, cfg)

	//连接失败，开始重连，
	if err != nil {
		connCount++
		if connCount > c.config.RetryMaxTries {
			Global.Logger["MQ"].Error("MQ建立连接失败：" + err.Error())
		}

		time.Sleep(c.config.RetryBaseTime * time.Millisecond)
		continue
	}

}
