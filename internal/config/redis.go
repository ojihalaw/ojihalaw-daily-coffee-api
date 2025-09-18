package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedisClient(config *viper.Viper, log *logrus.Logger) *redis.Client {
	addr := fmt.Sprintf("%s:%d", config.GetString("REDIS_HOST"), config.GetInt("REDIS_PORT"))

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		// Password: config.GetString("REDIS_PASSWORD"),
		DB: config.GetInt("REDIS_DB"),
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
		panic(err)
	}

	log.Info("Connected redis")
	return client

}
