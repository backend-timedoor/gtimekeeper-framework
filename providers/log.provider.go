package providers

import (
	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/log"
)

type LogServiceProvider struct{}

func (p *LogServiceProvider) Boot() {}

func (p *LogServiceProvider) Register() {
	app.Log = log.BootLog()
}
