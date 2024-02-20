package server

import (
	"github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
	"github.com/backend-timedoor/gtimekeeper-framework/base/server/servers"
	"github.com/backend-timedoor/gtimekeeper-framework/base/validation"
)

type Server struct {
	Servers    map[string]any
	Validation *validation.Validation
}

func (s *Server) Grpc() contracts.ServerHandle {
	return s.Servers["Grpc"].(*servers.Grpc)
}

func (s *Server) Http() contracts.ServerHandle {
	server := s.Servers["Http"].(*servers.Http)
	server.Server.Validator = s.Validation

	return server
}

func (s *Server) RegisterCustomeValidation(validations []contracts.CustomeValidation) {
	s.Validation.RegisterCustomeValidation(validations)
}
