package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"alligator/pkg/config"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func Init(opt config.Options) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     opt.Cache.Addr,
		Password: opt.Cache.Passwd,
		DB:       opt.Cache.DB,
	})
}

func Get(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func Set(key, value string, expire time.Duration) error {
	return rdb.Set(ctx, key, value, expire).Err()
}

func Del(key ...string) error {
	return rdb.Del(ctx, key...).Err()
}
