package Global

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Utils/RabbitMQ"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/token"
	"golang.org/x/time/rate"

	"AITranslatio/Utils/zipkin"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Infrastructure struct {
	Logger           map[string]*zap.Logger
	Config           interf.ConfigInterface
	DbConfig         interf.ConfigInterface
	DbClient         *gorm.DB
	RedisClient      *redis.Client
	Scripts          map[string]*redis.Script
	RabbitmqClient   *RabbitMQ.Client
	JwtManager       *token.JWTGenerator
	SnowflakeManager *SnowFlak.SnowFlakeGenerator
	Tracing          *zipkin.Tracing
	EncryptKey       []byte
	Limiter          *rate.Limiter
}

var (
	infra *Infrastructure
)

func InitInfrastructure(cfg interf.ConfigInterface, logger map[string]*zap.Logger, dbClient *gorm.DB, redisClient *redis.Client, scripts map[string]*redis.Script, MQClient *RabbitMQ.Client, tracing *zipkin.Tracing, l *rate.Limiter, one *sync.Once) *Infrastructure {

	one.Do(func() {
		infra = &Infrastructure{
			Config:           cfg,
			Logger:           logger,
			DbClient:         dbClient,
			RedisClient:      redisClient,
			RabbitmqClient:   MQClient,
			Scripts:          scripts,
			Tracing:          tracing,
			EncryptKey:       []byte(os.Getenv("PATHEXT")),
			SnowflakeManager: SnowFlak.CreateSnowflakeFactory(cfg, logger["Business"]),
			Limiter:          l,
		}

		AkOutTime := cfg.GetDuration("Token.AkOutTime")
		RkOutTime := cfg.GetDuration("Token.RkOutTime")

		c := &token.CreateToken{
			infra.EncryptKey,
			AkOutTime,
			RkOutTime,
			infra.SnowflakeManager,
			redisClient,
		}

		infra.JwtManager = token.CreateTokenFactory(c)

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
