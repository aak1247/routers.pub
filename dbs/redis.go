package dbs

import (
	"context"
	"github.com/go-redis/redis/v8"
	"routers.pub/env"
	"routers.pub/infra"
	"time"
)

const (
	Separator = ":" // key分割符号
)

var (
	RedisDb *redis.Client
	Redis   *env.Redis
)

func InitRedis() error {
	Redis = env.Conf.Redis
	address := Redis.Host + ":" + Redis.Port
	RedisDb = redis.NewClient(
		&redis.Options{
			Addr:         address,
			Password:     Redis.Password,
			PoolSize:     Redis.MaxActive, // 最大连接数
			MinIdleConns: 10,              // 最小空闲连接数
		},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisDb.Ping(ctx).Result()
	return err
}

func Get(key string) string {
	val, _ := RedisDb.Get(context.Background(), key).Result()
	return val
}
func Set(key string, value interface{}, expiration time.Duration) error {
	err := RedisDb.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		infra.Log.Error(err)
		return err
	}
	return nil
}
func Delete(key string) int64 {
	return RedisDb.Del(context.Background(), key).Val()
}

func Incr(key string) int64 {
	return RedisDb.Incr(context.Background(), key).Val()
}

func IncrBy(key string, incrBy int64) int64 {
	return RedisDb.IncrBy(context.Background(), key, incrBy).Val()
}

func Expire(key string, expiration time.Duration) bool {
	return RedisDb.Expire(context.Background(), key, expiration).Val()
}
