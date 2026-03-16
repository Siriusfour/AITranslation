package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() *redis.Client {

	Options := &redis.Options{
		Addr:     viper.GetString("Redis.DNS"),
		Password: viper.GetString("Redis.password"),
		DB:       viper.GetInt("Redis.DB"),
	}

	return redis.NewClient(Options)

}
