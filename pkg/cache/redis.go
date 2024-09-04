package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rc *redis.Client

func ConnectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     ":6381",
		Password: "",
		DB:       0,
	})
	rc = client
	fmt.Println("Connect redis success ...")
}

func SetRedis(ctx context.Context, key string, data interface{}, exp time.Duration) error {
	err := rc.Set(ctx, key, data, exp).Err()
	if err != nil {
		return errors.New("failed to set value redis")
	}

	return nil
}

func GetRedis(ctx context.Context, key string) (string, error) {
	val, err := rc.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("key does not exist")
	}
	if err != nil {
		return "", errors.New("failed to get value from redis")
	}

	return val, nil
}

func GenKeyRedis(prefix, val string) string {
	return fmt.Sprintf("%s_%s", prefix, val)
}
