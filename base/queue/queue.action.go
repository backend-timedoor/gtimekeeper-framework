package queue

import (
	"encoding/json"
	"time"

	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/hibiken/asynq"
)

type Queue struct {
	client *asynq.Client
	server *asynq.Server
	task *asynq.ServeMux
	Tasks []Task
}

type Task struct {
	Job contracts.Job
	Task *asynq.Task
}

func (q *Queue) Run() {
	// defer q.client.Close()

	for _, task := range q.Tasks {
		// init client
		options := append([]asynq.Option{
			asynq.MaxRetry(5),
			asynq.Timeout(3 * time.Minute),
		}, task.Job.Options()...)

		_, err := q.client.Enqueue(task.Task, options...)
		if err != nil {
			app.Log.Fatalf("could not enqueue task %s: %v", task.Job.Signature(), err)
		}

		// init handler 
		q.task.HandleFunc(task.Job.Signature(), task.Job.Handle)
	}

	if err := q.server.Run(q.task); err != nil {
		app.Log.Fatalf("could not run server: %v", err)
    }
}

func (q *Queue) Job(job contracts.Job, args any) {
	payload, err := json.Marshal(args)
	if err != nil {
		app.Log.Fatalf("queue %s marshal args error: %v", job.Signature(), err)
	}

	task := Task{
		Job:  job,
		Task: asynq.NewTask(job.Signature(), payload),
	}

	q.Tasks = append(q.Tasks, task)
}