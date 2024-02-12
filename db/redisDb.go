package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	redisAddr string
	client    *redis.Client
}

func NewRedisDB(redisAddr string) *RedisDB {
	return &RedisDB{
		redisAddr: redisAddr,
	}
}

func (r *RedisDB) Connect(ctx context.Context) error {
	r.client = redis.NewClient(&redis.Options{
		Addr: r.redisAddr,
	})
	return r.client.Ping(ctx).Err()
}

func (r *RedisDB) GetDB() *redis.Client {
	return r.client
}

func (r *RedisDB) Close() {
	_ = r.client.Close()
}
