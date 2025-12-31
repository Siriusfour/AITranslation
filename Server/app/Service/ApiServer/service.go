package ApiServer

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Utils/RabbitMQ"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/app/DAO/ApiDAO"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ApiServer struct {
	cfg                interf.ConfigInterface
	logger             *zap.Logger
	redis              *redis.Client
	RabbitMQ           *RabbitMQ.Client
	scripts            map[string]*redis.Script
	snowFlakeGenerator *SnowFlak.SnowFlakeGenerator
	DAO                *ApiDAO.ApiDAO
}

func NewService(cfg interf.ConfigInterface, logger *zap.Logger, redisClient *redis.Client, rabbit *RabbitMQ.Client, SnowFlakeGenerator *SnowFlak.SnowFlakeGenerator, scripts map[string]*redis.Script, DAO *ApiDAO.ApiDAO) *ApiServer {
	return &ApiServer{
		logger:             logger,
		cfg:                cfg,
		redis:              redisClient,
		RabbitMQ:           rabbit,
		scripts:            scripts,
		snowFlakeGenerator: SnowFlakeGenerator,
		DAO:                DAO,
	}
}
