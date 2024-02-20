package redis

import (
	baseRedis "github.com/go-redis/redis"
)

type Config struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}

func New(config *Config) (*Redis, error) {
	cache := baseRedis.NewClient(&baseRedis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := cache.Ping().Result()
	if err != nil {
		return nil, err
	}

	prefix := config.Prefix
	if prefix == "" {
		prefix = "gtimekeeper"
	}

	return &Redis{
		Redis:  cache,
		Prefix: prefix,
	}, nil
}
