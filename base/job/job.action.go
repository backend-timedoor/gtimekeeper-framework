package job

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database/redis"
	"github.com/backend-timedoor/gtimekeeper-framework/base/job/custom"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/hibiken/asynq"
)

const CACHE_KEY = "job-module"

type Config struct {
	// RedisOpt    *asynq.RedisClientOpt
	ServerOpt   *asynq.Config
	ScheduleOpt *asynq.SchedulerOpts
}

type Job struct {
	cache     *redis.Redis
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
	mux       *asynq.ServeMux
}

func (j *Job) RegisterQueue(queues []contracts.Queue) {
	queues = append(queues, []contracts.Queue{
		&custom.EmailJob{},
	}...)

	for _, queue := range queues {
		// check signature
		if j.containSignature(queue.Signature()) {
			log.Fatalf("job with signature %s is already exists", queue.Signature())
		}

		// register queue
		err := j.cache.Push(CACHE_KEY, queue.Signature())
		if err != nil {
			log.Fatalf("cannot register queue %s: %v", queue.Signature(), err)
		}
	}
}

func (j *Job) RegisterSchedule(schedules []contracts.Schedule) {
	for _, schedule := range schedules {
		// check signature
		if j.containSignature(schedule.Signature()) {
			log.Fatalf("job with signature %s is already exists", schedule.Signature())
		}

		// register scheduler
		_, err := j.scheduler.Register(
			schedule.Schedule(),
			asynq.NewTask(schedule.Signature(), nil),
			schedule.Options()...,
		)

		j.mux.HandleFunc(schedule.Signature(), schedule.Handle)

		if err != nil {
			log.Fatalf("cannot register schedule %s: %v", schedule.Signature(), err)
		}

		err = j.cache.Push(CACHE_KEY, schedule.Signature())
		if err != nil {
			log.Fatalf("cannot register schedule %s: %v", schedule.Signature(), err)
		}
	}
}

func (j *Job) Queue(job contracts.Queue, args any) error {
	if !j.containSignature(job.Signature()) {
		return fmt.Errorf("job with signature %s is unregistered", job.Signature())
	}

	payload, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("queue %s marshal args error: %v", job.Signature(), err)
	}

	options := append([]asynq.Option{
		asynq.MaxRetry(5),
		asynq.Timeout(2 * time.Minute),
		// asynq.ProcessIn(time.Second),
	}, job.Options()...)

	_, err = j.client.Enqueue(asynq.NewTask(job.Signature(), payload), options...)
	if err != nil {
		return fmt.Errorf("queue %s client error enqueue error: %v", job.Signature(), err)
	}

	j.mux.HandleFunc(job.Signature(), job.Handle)
	j.cache.Push(CACHE_KEY, job.Signature())
	container.Log().Infof("task added %s", job.Signature())

	return nil
}

func (j *Job) Run() error {
	if err := j.scheduler.Start(); err != nil {
		return fmt.Errorf("could not run scheduler: %v", err)
	}

	if err := j.server.Start(j.mux); err != nil {
		return fmt.Errorf("could not run server: %v", err)
	}

	return nil
}

func (j *Job) containSignature(signature string) bool {
	jobs := j.cache.Retrieve(CACHE_KEY)

	for _, j := range jobs {
		if j == signature {
			return true
		}
	}

	return false
}
