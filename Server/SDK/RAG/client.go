package RAG

import (
	PB "AITranslatio/SDK/RAG/pb"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
	"time"
)

type Client struct {
	conn   *grpc.ClientConn
	client PB.AskServiceClient
}

type Option func(*ClientOptions)

type ClientOptions struct {
	address     string
	timeout     time.Duration
	dialOptions []grpc.DialOption
}

// 三个配置函数
func WithAddress(addr string) Option {
	return func(o *ClientOptions) {
		if strings.TrimSpace(addr) != "" {
			o.address = addr
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *ClientOptions) {
		if timeout > 0 {
			o.timeout = timeout
		}
	}
}

func WithDialOption(opt grpc.DialOption) Option {
	return func(o *ClientOptions) {
		if opt != nil {
			o.dialOptions = append(o.dialOptions, opt)
		}
	}
}

// 新建客户端
func NewClient(opts ...Option) (*Client, error) {
	options := &ClientOptions{
		address: "127.0.0.1:7071",
		timeout: 15 * time.Second,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(options)
		}
	}
	if len(options.dialOptions) == 0 {
		options.dialOptions = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), options.timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, options.address, options.dialOptions...)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: PB.NewAskServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *Client) Ask(ctx context.Context, question string, contextID string) (string, error) {

	if c == nil || c.client == nil {
		return "", errors.New("client is nil")
	}
	if strings.TrimSpace(question) == "" {
		return "", errors.New("question is required")
	}
	req := &PB.AskRequest{
		Question:  question,
		SessionId: contextID,
	}
	resp, err := c.client.Ask(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.GetAnswer(), nil
}
