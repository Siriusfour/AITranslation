package Global

import (
	"AITranslatio/Utils/SSE"
	"AITranslatio/config/interf"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

var Logger *zap.Logger
var MySQL_Client *gorm.DB
var PostgreSQL_Client *gorm.DB
var SQLserver_Client *gorm.DB
var RedisClient *redis.Client
var PKEY = []byte(os.Getenv("PATHEXT"))
var SSEClients *SSE.SSEClients
var Config interf.ConfigInterface
var DB_Config interf.ConfigInterface
var DataFormt = "2006-01-02 15:04:05"
