package providers

import (
	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/database"
)

type DatabaseServiceProvider struct{}

func (p *DatabaseServiceProvider) Boot() {}

func (p *DatabaseServiceProvider) Register() {
	app.DB = database.BootDatabase()
}
