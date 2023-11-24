package providers

import (
	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/backend-timedoor/gtimekeeper/base/log"
)

type LogServiceProvider struct{}

func (p *LogServiceProvider) Boot() {}

func (p *LogServiceProvider) Register() {
	app.Log = log.BootLog()
}
