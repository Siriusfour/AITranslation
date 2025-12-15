package RabbitMQ

import (
	"AITranslatio/Global"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

// mock connection 对象
type mockConn struct {
	closed bool
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

// mock 成功的 Dial
func mockDialSuccess(uri string, cfg amqp.Config) (*amqp.Connection, error) {
	return &amqp.Connection{}, nil
}

// mock 失败的 Dial
func mockDialFail(uri string, cfg amqp.Config) (*amqp.Connection, error) {
	return nil, errors.New("mock dial failed")
}

// mock 发送消息成功，没有ack
func mockSendMessageSuccess(c *Client) error {

}

func mockSendMessageFail(c *Client) error {

}

func initTestClient(DialMQFunc DialMQ) *Client {

	if Global.Logger == nil {
		Global.Logger = make(map[string]*zap.Logger)
	}
	Global.Logger["MQ"] = zap.NewNop() // ✅ 这个 Logger 不打印，也不会 panic

	config := &Config{
		URI:           "amqp://@127.0.0.1:5672/",
		Durable:       false,
		Workers:       10,
		RetryMaxTries: 10,
		RetryWaitTime: 10000 * time.Millisecond,
		RetryBaseTime: 100 * time.Millisecond,
		EnableConfirm: true,
	}

	client := &Client{
		Wg:     sync.WaitGroup{},
		Config: *config,
		Conn:   nil,
		Closed: atomic.Bool{},
		DialMQ: DialMQFunc,
	}

	return client
}

// 测试成功连接默认值逻辑是否正确
//func TestClientConnect_DefaultValues(t *testing.T) {
//
//	c := initTestClient(mockDialSuccess)
//
//	err := c.Connect()
//	if err != nil {
//		t.Fatalf("没有正确连接： %v", err)
//	}
//
//	// 默认值验证
//	if c.Config.Heartbeat <= 0 {
//		t.Errorf("初始心跳不正确，Heartbeat：%d", c.Config.Heartbeat)
//	}
//}

// mock 不停失败，用于重试测试
//func TestClientConnect_RetryExceeded(t *testing.T) {
//
//	c := initTestClient(mockDialFail)
//
//	start := time.Now()
//	err := c.Connect()
//	elapsed := time.Since(start)
//
//	//连接成功则返回错误
//	if err == nil {
//		t.Fatalf("expected error, got nil")
//	}
//
//	if elapsed < 20*time.Millisecond {
//		t.Errorf("expected delay (backoff), got too fast: %v", elapsed)
//	}
//}

// 测试重连功能
//func TestClientReconnect(t *testing.T) {
//
//	c := initTestClient(mockDialSuccess)
//
//	err := c.Connect()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	//延迟10s，断开连接
//	time.Sleep(1 * time.Second)
//	c.notifychannel <- &amqp.Error{Reason: "测试重连功能是否正常"}
//
//	if c.Conn == nil {
//		t.Fatal("没有主动重连")
//	}
//	//TODO 调用MQ生产，消费逻辑,看是否正确
//}
