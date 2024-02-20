package job

import (
	"log"

	"github.com/backend-timedoor/gtimekeeper-framework/base/database"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/hibiken/asynq"
)

const ContainerName string = "job"

func New(config *Config) *Job {
	db := container.App["db"].(*database.Database)
	if db.Redis == nil {
		log.Fatal("job module need redis, Redis is not initialized")
	}

	dbRedisOpts := db.Config.Redis
	redisOpts := &asynq.RedisClientOpt{
		Addr:     dbRedisOpts.Addr,
		Password: dbRedisOpts.Password,
		DB:       dbRedisOpts.DB,
	}

	j := &Job{
		mux:    asynq.NewServeMux(),
		client: asynq.NewClient(redisOpts),
		server: asynq.NewServer(
			redisOpts,
			*config.ServerOpt,
		),
		scheduler: asynq.NewScheduler(
			redisOpts,
			config.ScheduleOpt,
		),
		cache: db.Redis,
	}

	j.cache.Forget(CACHE_KEY)

	container.Set(ContainerName, j)

	return j
}
