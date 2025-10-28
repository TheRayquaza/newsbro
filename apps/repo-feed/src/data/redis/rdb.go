package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Sentinels    []string // ["localhost:26379", "localhost:26380", "localhost:26381"]
	MasterName   string   // "mymaster"
	Password     string
	DB           int // 0
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Initialize(config RedisConfig) (*redis.Client, error) {
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    config.MasterName,
		SentinelAddrs: config.Sentinels,
		Password:      config.Password,
		DB:            config.DB,
		DialTimeout:   config.DialTimeout,
		ReadTimeout:   config.ReadTimeout,
		WriteTimeout:  config.WriteTimeout,

		SentinelPassword: config.Password,

		PoolSize:     10,
		MinIdleConns: 2,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

func HealthCheck(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}
