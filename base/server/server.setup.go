package server

import (
	"reflect"

	"github.com/backend-timedoor/gtimekeeper/base/contracts"
)


func BootServer(servers []contracts.ServerHandle) contracts.Server {
	s := &Server{}
	serverBag := map[string]any{}

	for _, server := range servers {
		server.Start()
		server.Handler()

		serverBag[reflect.TypeOf(server).Elem().Name()] = server
	}

	s.servers = serverBag

	return s
}