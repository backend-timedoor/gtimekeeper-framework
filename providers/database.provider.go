package providers

import (
	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/backend-timedoor/gtimekeeper/base/database"
)

type DatabaseServiceProvider struct{}

func (p *DatabaseServiceProvider) Boot() {}

func (p *DatabaseServiceProvider) Register() {
	app.DB = database.BootDatabase()
}
