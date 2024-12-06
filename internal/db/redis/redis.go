package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis_rate/v10"

	"github.com/redis/go-redis/v9"
	logger "github.com/sirupsen/logrus"

	"one-day-server/configs"
)

const ErrorEqualOrSmallerStreamID = "ERR The ID specified in XADD is equal or smaller than the target stream top item"

var NilError = errors.New("key does not exist, err: redis: nil")

var (
	client *RedisClient
	once   sync.Once
)

func GetClient() IClient {
	once.Do(func() {
		rdb := connectRedis()
		client = &RedisClient{
			rdb:         rdb,
			RateLimiter: redis_rate.NewLimiter(rdb),
		}
	})
	return client
}

type IClient interface {
	// redis KV
	GetResult(ctx context.Context, key string) (string, error)
	WriteResult(ctx context.Context, key, value string) error
	WriteResultWithTTL(ctx context.Context, key, value string, ttl time.Duration) error
	GetRateLimiter() *redis_rate.Limiter

	// no need to wrap all the methods, allow expose client to outer package
	GetRDB() *redis.Client
}

type RedisClient struct {
	rdb         *redis.Client
	RateLimiter *redis_rate.Limiter
}

func connectRedis() *redis.Client {
	dbInt, err := strconv.Atoi(configs.GetEnvDefault("REDIS_DB", "0"))
	if err != nil {
		logger.Warnf("Fail to parse redis db: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.MustGetEnv("REDIS_HOST"), configs.MustGetEnv("REDIS_PORT")),
		Username: configs.MustGetEnv("REDIS_USERNAME"),
		Password: configs.MustGetEnv("REDIS_PASSWORD"),
		DB:       dbInt,
	})

	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Panicf("connect redis failed: %s", err)
	}
	logger.Infof("connect to Redis: %s\n", pong)
	return rdb
}

func (rc *RedisClient) GetResult(ctx context.Context, key string) (string, error) {
	result, err := rc.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", NilError
	} else if err != nil {
		return "", fmt.Errorf("read result failed, err: %s", err)
	} else {
		return result, nil
	}
}

func (rc *RedisClient) WriteResult(ctx context.Context, key, value string) error {
	retry := 3
	var err error
	for retry > 0 {
		err = rc.rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			logger.Errorf("write to redis %s failed, value: %s, err: %s", key, value, err)
			retry--
			continue
		}
		break
	}
	return err
}

func (rc *RedisClient) WriteResultWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	retry := 3
	var err error
	for retry > 0 {
		err = rc.rdb.Set(ctx, key, value, ttl).Err()
		if err != nil {
			logger.Errorf("write to redis %s failed, value: %s, err: %s", key, value, err)
			retry--
			continue
		}
		break
	}
	return err
}

func (rc *RedisClient) GetRateLimiter() *redis_rate.Limiter {
	return rc.RateLimiter
}

func (rc *RedisClient) GetRDB() *redis.Client {
	return rc.rdb
}
