package boot

import (
	"reflect"

	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/providers"
)

func Booting(pvds []contracts.ServiceProvider) {
	pvds = append(pvds, []contracts.ServiceProvider{
		&providers.CacheServiceProvider{},
		&providers.DatabaseServiceProvider{},
		&providers.MailServiceProvider{},
	}...)

	for _, provider := range pvds {
		provider.Register()
		r := reflect.TypeOf(provider).Elem()
		if (r.Name() == "ConfigServiceProvider") {
			log := providers.LogServiceProvider{}
			log.Register()
			log.Boot()
		}

		provider.Boot()
	}
}
