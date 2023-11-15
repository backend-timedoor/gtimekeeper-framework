package cache

import (
	"fmt"

	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/go-redis/redis"
)

func BootCache(prefix string) contracts.Cache {
	config := app.Config.Get("database.redis").(map[string]any)

	cache := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config["host"].(string), config["port"].(string)),
		Password: config["password"].(string),
		DB: 0,
	})

	_, err := cache.Ping().Result()
	if err != nil {
		// err := fmt.Errorf("failed to link redis:%s\n%+v", err, string(debug.Stack()))
		// fmt.Println(err.Error())
		app.Log.Errorf("failed to link redis: %s", err)
	}

	if prefix == "" {
		prefix = "gtime_keeper"
	}

	return &Redis{
		Redis: cache,
		Prefix: prefix,
	}
}
