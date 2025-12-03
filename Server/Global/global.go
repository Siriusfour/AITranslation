package Global

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Utils/RabbitMQ"

	"AITranslatio/Utils/zipkin"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Infrastructure struct {
	Logger map[string]*zap.Logger

	Config   interf.ConfigInterface
	DbConfig interf.ConfigInterface

	DbClient *gorm.DB

	RedisClient    *redis.Client
	RabbitmqClient *RabbitMQ.Client

	//snowflakeManage *SnowFlak.SnowFlake
	Tracing *zipkin.Tracing

	EncryptKey []byte
}

var (
	infra     *Infrastructure
	infraOnce sync.Once
)

func InitInfrastructure(cfg interf.ConfigInterface, logger map[string]*zap.Logger, dbClient *gorm.DB, redisClient *redis.Client, MQClient *RabbitMQ.Client, tracing *zipkin.Tracing) *Infrastructure {

	infraOnce.Do(func() {
		infra = &Infrastructure{
			Config:         cfg,
			Logger:         logger,
			DbClient:       dbClient,
			RedisClient:    redisClient,
			RabbitmqClient: MQClient,
			Tracing:        tracing,
			EncryptKey:     []byte(os.Getenv("PATHEXT")),
		}
	})

	return infra
}

func GetInfra() *Infrastructure {
	return infra
}

func (infra *Infrastructure) GetConfig() interf.ConfigInterface {
	return infra.Config
}

func (infra *Infrastructure) GetLogger(name string) *zap.Logger {
	return infra.Logger[name]
}
