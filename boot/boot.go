package boot

import (
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
)

func Booting(pvds []contracts.ServiceProvider) {
	container.App = map[string]any{}
	pvds = append(pvds, []contracts.ServiceProvider{
		// &providers.CacheServiceProvider{},
		// &providers.DatabaseServiceProvider{},
		// &providers.MailServiceProvider{},
	}...)

	for _, provider := range pvds {
		provider.Register()
		provider.Boot()
	}
}
