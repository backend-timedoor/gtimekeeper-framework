package custom

import (
	"context"
	"reflect"

	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/hibiken/asynq"
)

type EmailJob struct{}

func (m *EmailJob) Signature() string {
	return "email:job"
}

func (m *EmailJob) Options() []asynq.Option {
	return []asynq.Option{}
}

func (m *EmailJob) Handle(ctx context.Context, t *asynq.Task) error {
	log := container.Log()
	mail := container.ExecRef("mail", "SendWithQueue", []reflect.Value{
		reflect.ValueOf(t.Payload()),
	})

	err := mail[0].Interface()
	if err != nil {
		log.Errorf("mail queue error %v", err.(error).Error())
	}

	return nil
}
