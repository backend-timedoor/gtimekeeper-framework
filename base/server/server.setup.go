package server

import (
	"reflect"

	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/backend-timedoor/gtimekeeper/base/server/validation"
)


func BootServer(servers []contracts.ServerHandle) contracts.Server {
	s := &Server{}
	serverBag := map[string]any{}

	for _, server := range servers {
		server.Start()
		server.Handler()

		serverBag[reflect.TypeOf(server).Elem().Name()] = server
	}

	s.Servers = serverBag
	s.Validation = validation.BootCustomValidation()

	return s
}