package Global

import (
	"AITranslatio/Utils/UtilsStruct"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"time"
)

var Logger *zap.SugaredLogger
var OutTime = 30 * 60 * time.Second
var DB *gorm.DB
var RedisClient *redis.Client
var PKEY = []byte(os.Getenv("PATHEXT"))
var TokenMap *UtilsStruct.TokenMap
