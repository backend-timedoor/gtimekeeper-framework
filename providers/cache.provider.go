package providers

import (
	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/cache"
)

type CacheServiceProvider struct{}
func (p *CacheServiceProvider) Boot() {}

func (p *CacheServiceProvider) Register() {
	app.Cache = cache.BootCache()
}
