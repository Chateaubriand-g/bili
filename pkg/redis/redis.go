package redis

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

var (
	redisC *redis.Client
	ctx = context.Background()
	//once sync.Once
	initialized bool
)

var (
	dbIndex = 0
	poolSize = 10
	dialTimeout = 5
	readTimeout = 3
	writeTimeout = 3
)

func InitRedis(cfg *Config) {
	redisC := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		PassWord: cfg.RedisPassword,
		DB: dbIndex,
		PoolSize: poolSize,
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
	})

	if err:= redisC.Ping(ctx).Err();err!=nil{
		log.Fatal("redie connect failed: %w",error)
	}
	initialized = true
}

func Set(ctx context.Context,key string,value interface{},ttl time.Duration) error {
	if !initialized{
		return redis.Nil
	}
	return redisC.Set(ctx,key,val,ttl).Err()
}

func Get(ctx context.Context,key string) (string,error) {
	if !initialized{
		return "",redis.Nil
	}
	return redisC.Get(ctx,key).Result()
}

func Del(ctx context.Context,keys ...string) error {
	if !initialized{
		return "",redis.Nil
	}
	return redisC.Del(ctx,keys...).Err()
}

func Exists(ctx context.Context,keys ...string) (int64,error) {
	if !initialized {
		return 0,redis.Nil
	}
	return redisc.Exists(ctx,keys...).Result()
}

func Ince(ctx context.Context,key string) (int64,error) {
	if !initialized{
		return 0,redis.Nil
	}
	return rdb.Incr(ctx,key).Result
}

func Expire(ctx context.Context,key string,ttl time.Duration) error {
	if !initialized{
		return redis.Nil
	}
	return redisC.Expire(ctx,key,ttl).Err()
}