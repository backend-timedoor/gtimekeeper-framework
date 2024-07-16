package exceptions

import (
	"encoding/json"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper/types/protocol"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
)

type ErrorMessage struct {
	StatusCode any    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
	Errors     any    `json:"errors,omitempty"`
}

type GTimeError struct {
	Protocol protocol.Protocol `json:"-"`
	HttpCode int               `json:"-"`
	GrpcCode codes.Code        `json:"-"`
	message  *ErrorMessage
}

func (e *GTimeError) Error() string {
	err, _ := json.Marshal(e.message)

	return string(err)
}

func (e *GTimeError) SetCodes(code any) {
	var (
		httpCode int
		grpcCode codes.Code
	)

	ref := reflect.TypeOf(code).Kind()
	switch ref {
	case reflect.Int:
		httpCode = code.(int)
		grpcCode = httpToGrpcStatus[httpCode]
	default:
		grpcCode = code.(codes.Code)
		httpCode = grpcToHttpStatus[grpcCode]
	}

	e.HttpCode = httpCode
	e.GrpcCode = grpcCode
}

func (e *GTimeError) Make(code any, message *ErrorMessage) *GTimeError {
	e.SetCodes(code)
	e.message = message

	return e
}

func (e *GTimeError) simpleReturn(ptc protocol.Protocol) error {
	if ptc == protocol.HTTP {
		return e.HttpError()
	} else {
		return e.GrpcError()
	}
}

func (e *GTimeError) HttpError() error {
	return echo.NewHTTPError(e.HttpCode, e.message)
}

func (e *GTimeError) GrpcError() error {
	return status.Error(e.GrpcCode, e.Error())
}
