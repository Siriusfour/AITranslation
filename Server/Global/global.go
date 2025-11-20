package Global

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Utils/SSE"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

var Logger map[string]*zap.Logger
var MySQL_Client *gorm.DB
var PostgreSQL_Client *gorm.DB

var RedisClient *redis.Client
var EncryptKey = []byte(os.Getenv("PATHEXT"))
var SSEClients *SSE.SSEClients
var Config interf.ConfigInterface
var DB_Config interf.ConfigInterface
var RabbitmqClient any

var DataFormt = "2006-01-02 15:04:05"
