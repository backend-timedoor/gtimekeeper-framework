package queue

import (
	"fmt"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/hibiken/asynq"
)

func BootQueue() contracts.Queue {
	var queue Queue
	config := app.Config.Get("database.redis").(map[string]any)
	redis := asynq.RedisClientOpt{
        Addr: fmt.Sprintf("%s:%s", config["host"].(string), config["port"].(string)),
		Password: config["password"].(string),
		DB: app.Config.GetInt("database.redis.database", 0),
    }
	
	// init client
	queue.client = asynq.NewClient(redis)

	// init server
	queue.server = asynq.NewServer(
        redis,
        asynq.Config{},
    )

	// init mux server
	queue.task = asynq.NewServeMux()

	return &queue
}