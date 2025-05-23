package Config

import (
	"AITranslatio/Global"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() {

	Options := &redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	Global.RedisClient = redis.NewClient(Options)

}
