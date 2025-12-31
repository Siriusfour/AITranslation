package bootstrap

import (
	"AITranslatio/app/Service/ApiServer"
	"context"
)

func InitConsumer(ctx context.Context, s *ApiServer.ApiServer) {

	ch := make(chan error)
	s.RabbitMQ.Consumer(ctx, "seckill_order_queue", s.SeckillOrderHandler, ch)
}
