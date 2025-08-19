package Global

import (
	"AITranslatio/Utils/SSE"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

var Logger *zap.SugaredLogger
var DB *gorm.DB
var RedisClient *redis.Client
var PKEY = []byte(os.Getenv("PATHEXT"))
var SSEClients *SSE.SSEClients
var DataFormt = "2006-01-02 15:04:05"
