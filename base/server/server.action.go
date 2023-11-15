package server

import (
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/backend-timedoor/gtimekeeper/base/server/servers"
)

type Server struct {
	servers map[string]any
}

func (s *Server) Grpc() contracts.ServerHandle {
	return s.servers["Grpc"].(*servers.Grpc)
}

func (s *Server) Http() contracts.ServerHandle {
	return s.servers["Http"].(*servers.Http)
}