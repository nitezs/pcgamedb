package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nitezs/pcgamedb/config"
	"github.com/nitezs/pcgamedb/log"

	"github.com/redis/go-redis/v9"
)

var cache *redis.Client
var mutx = &sync.RWMutex{}

func connect() {
	if !config.Config.RedisAvaliable {
		return
	}
	cache = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DBIndex,
	})
	err := HealthCheck()
	if err != nil {
		log.Logger.Panic("Cannot connect to redis")
	}
	log.Logger.Info("Connected to redis")
}

func CheckConnect() {
	mutx.RLock()
	if cache != nil {
		mutx.RUnlock()
		return
	}
	mutx.RUnlock()

	mutx.Lock()
	if cache == nil {
		connect()
	}
	mutx.Unlock()
}

func HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := cache.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if result != "PONG" {
		return fmt.Errorf("unexpected response from Redis: %s", result)
	}
	return nil
}

func Get(key string) (string, bool) {
	CheckConnect()
	ctx := context.Background()
	value, err := cache.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return value, true
}

func Add(key string, value interface{}) error {
	CheckConnect()
	ctx := context.Background()
	cmd := cache.Set(ctx, key, value, 7*24*time.Hour)
	return cmd.Err()
}

func AddWithExpire(key string, value interface{}, expire time.Duration) error {
	CheckConnect()
	ctx := context.Background()
	cmd := cache.Set(ctx, key, value, expire)
	return cmd.Err()
}
