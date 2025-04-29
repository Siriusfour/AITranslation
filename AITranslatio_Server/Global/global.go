package Global

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

var Logger *zap.SugaredLogger
var OutTime = 30 * 60 * time.Second
var DB *gorm.DB
var RedisClient *redis.Client
