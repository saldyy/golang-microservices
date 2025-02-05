package database

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisInstance *RedisClientInstance

type RedisClientInstance struct {
	Client *redis.Client
  logger *slog.Logger
}

type RedisConfigOptions struct {
	options redis.Options
}

func (opt RedisConfigOptions) WithAddr(addr string) RedisConfigOptions {
	opt.options.Addr = addr
	return opt
}

func (opt RedisConfigOptions) WithPassword(password string) RedisConfigOptions {
	opt.options.Password = password
	return opt
}

func (opt RedisConfigOptions) WithDB(db int) RedisConfigOptions {
	opt.options.DB = db
	return opt
}

func (opt RedisConfigOptions) GetConfig() redis.Options {
	return opt.options
}

func NewRedisConfigOptions() RedisConfigOptions {
	redisOpts := redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	conf := RedisConfigOptions{options: redisOpts}

	return conf
}

func NewRedisClientInstance(logger * slog.Logger) *RedisClientInstance {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPass := os.Getenv("REDIS_PASSWORD")

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	config := NewRedisConfigOptions().WithAddr(redisAddr).WithPassword(redisPass).WithDB(redisDB).GetConfig()
	rdb := redis.NewClient(&config)

  logger.Info("Connected to Redis")
	return &RedisClientInstance{Client: rdb, logger: logger}
}

func (rci *RedisClientInstance) RedisHealth() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rci.Client.Ping(ctx).Result()
  rci.logger.Info("Redis ping", pong)

	if err != nil {
    rci.logger.Error("Error connecting to Redis", err)
		return map[string]string{
			"message": "not_ok",
		}
	}

	return map[string]string{
		"message": "ok",
	}
}
