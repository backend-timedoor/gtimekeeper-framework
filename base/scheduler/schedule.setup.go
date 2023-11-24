package schedule

import (
	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/robfig/cron/v3"
)


func BootSchedule(schedules []contracts.ScheduleEvent) contracts.Schedule {
	cron := cron.New()
	for _, schedule := range schedules {
		_, err := cron.AddFunc(schedule.Schedule(), schedule.Handle)
		if err != nil {
			app.Log.Fatalf("cannot register schedule %s: %v", schedule.Signature(), err)
		}

		app.Log.Infof("schedule %s registered", schedule.Signature())
	}

	return &Schedule{cron}
}

