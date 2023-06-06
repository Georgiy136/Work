package redis

import (
	"context"
	"fmt"
	"myapp/config"

	"github.com/redis/go-redis/v9"
)


func New(cfg config.Redis) (*redis.Client, error) {

	Rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if _, err := Rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("redis - NewRedis - Ping: %w", err)
	}

	return Rdb, nil
}
