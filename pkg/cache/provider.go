package cache

import (
	"context"
	"log"
	"service/config"
	"service/pkg/datastore/redis"
)

func NewCache(ctx context.Context, cfg *config.Config) ICache {
	if cfg.Cache.Driver == "redis" {

		return redis.NewRedis(ctx, cfg)
	}

	log.Fatal("not support cache driver")
	return nil
}
