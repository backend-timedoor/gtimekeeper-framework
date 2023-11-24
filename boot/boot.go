package boot

import (
	"reflect"

	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/backend-timedoor/gtimekeeper/providers"
)

func Booting(pvds []contracts.ServiceProvider) {
	pvds = append(pvds, []contracts.ServiceProvider{
		&providers.CacheServiceProvider{},
		&providers.DatabaseServiceProvider{},
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
