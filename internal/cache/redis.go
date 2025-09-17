package cache

import (
	"context"
	"errors"
	"time"

	"github.com/xuanviet96/seta-training/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func New(cfg config.Config, log *zap.Logger) (*redis.Client, error) {
	if cfg.RedisAddr == "" {
		return nil, errors.New("REDIS_ADDR empty")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	log.Info("connected to redis", zap.String("addr", cfg.RedisAddr))
	return rdb, nil
}

func TTL(cfg config.Config) time.Duration {
	ttl := time.Duration(cfg.RedisTTLSeconds) * time.Second
	if ttl <= 0 {
		return 5 * time.Minute
	}
	return ttl
}
