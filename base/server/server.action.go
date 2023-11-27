package server

import (
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/backend-timedoor/gtimekeeper/base/server/servers"
	"github.com/backend-timedoor/gtimekeeper/base/server/validation"
	"github.com/backend-timedoor/gtimekeeper/base/server/validation/custom"
)

type Server struct {
	Servers map[string]any
	Validation *validation.CustomeValidation
}

func (s *Server) Grpc() contracts.ServerHandle {
	return s.Servers["Grpc"].(*servers.Grpc)
}

func (s *Server) Http() contracts.ServerHandle {
	server := s.Servers["Http"].(*servers.Http)
	server.Server.Validator = s.Validation

	return server
}

func (s *Server) RegisterCustomeValidation(validations []contracts.Validation) {
	validations	= append(validations, []contracts.Validation{
		&custom.UniqueValidator{},
	}...)

	for _, validation := range validations {
		s.Validation.Validator.RegisterValidation(validation.Signature(), validation.Handle)
	}
}