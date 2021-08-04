package redisx

import (
	"context"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func NewRedis(c *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         c.Url(),
		DB:           c.Database,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.MinIdleConn,
	})
	log.Info("Connected to Redis")
	if err := client.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}
	return client
}
