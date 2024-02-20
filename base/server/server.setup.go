package server

import (
	"reflect"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/container"
)

const ContainerName string = "server"

func New(servers []contracts.ServerHandle) contracts.Server {
	s := &Server{}
	serverBag := map[string]any{}

	for _, server := range servers {
		server.Start()
		server.Handler()

		serverBag[reflect.TypeOf(server).Elem().Name()] = server
	}

	s.Servers = serverBag
	s.Validation = app.Validation

	container.Set(ContainerName, s)

	return s
}
