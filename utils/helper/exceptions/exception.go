package exceptions

import (
    "github.com/backend-timedoor/gtimekeeper-framework/utils/helper/types/protocol"
    "google.golang.org/grpc/codes"
    "net/http"
)

// Generic function to handle different errors
func NewError(code any, defaultMsg string, args ...any) error {
	ptc := protocol.GRPC
	msg := defaultMsg
	if len(args) > 0 {
		if msgArg, ok := args[0].(string); ok {
			msg = msgArg
		}

		if len(args) > 1 {
			if ptcArg, ok := args[1].(protocol.Protocol); ok {
				ptc = ptcArg
			}
		}
	}

	gTimeError := GTimeError{
		Protocol: ptc,
	}
	err := gTimeError.Make(code, &ErrorMessage{
		Message: msg,
	})

	return err.simpleReturn(ptc)
}

func NotFound(args ...any) error {
	return NewError(codes.NotFound, "not found", args...)
}

func InternalServerError(args ...any) error {
	return NewError(codes.Internal, "internal server error", args...)
}

func BadRequest(args ...any) error {
	return NewError(codes.InvalidArgument, "bad request", args...)
}

func Unauthorized(args ...any) error {
	return NewError(codes.Unauthenticated, "unauthorized", args...)
}

func Forbidden(args ...any) error {
	return NewError(codes.PermissionDenied, "forbidden", args...)
}

func Conflict(args ...any) error {
	return NewError(codes.Aborted, "conflict", args...)
}

func Gone(args ...any) error {
	return NewError(codes.Canceled, "gone", args...)
}

func TooManyRequests(args ...any) error {
	return NewError(codes.ResourceExhausted, "too many requests", args...)
}

func NotImplemented(args ...any) error {
	return NewError(codes.Unimplemented, "not implemented", args...)
}

func BadGateway(args ...any) error {
	return NewError(codes.Unavailable, "bad gateway", args...)
}

func ServiceUnavailable(args ...any) error {
	return NewError(codes.Unavailable, "service unavailable", args...)
}

func GatewayTimeout(args ...any) error {
	return NewError(codes.DeadlineExceeded, "gateway timeout", args...)
}

func HTTPVersionNotSupported(args ...any) error {
	return NewError(codes.Unimplemented, "HTTP version not supported", args...)
}

func RequestEntityTooLarge(args ...any) error {
	return NewError(codes.InvalidArgument, "request entity too large", args...)
}

func UnsupportedMediaType(args ...any) error {
	return NewError(codes.InvalidArgument, "unsupported media type", args...)
}

func UnprocessableEntity(args ...any) error {
	return NewError(codes.InvalidArgument, "unprocessable entity", args...)
}

func MethodNotAllowed(args ...any) error {
	return NewError(codes.PermissionDenied, "method not allowed", args...)
}

func PreconditionFailed(args ...any) error {
	return NewError(codes.FailedPrecondition, "precondition failed", args...)
}

var grpcToHttpStatus = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusRequestTimeout,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusRequestTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
}

var httpToGrpcStatus = map[int]codes.Code{
	http.StatusBadRequest:              codes.InvalidArgument,
	http.StatusUnauthorized:            codes.Unauthenticated,
	http.StatusForbidden:               codes.PermissionDenied,
	http.StatusNotFound:                codes.NotFound,
	http.StatusConflict:                codes.AlreadyExists,
	http.StatusGone:                    codes.Aborted,
	http.StatusTooManyRequests:         codes.ResourceExhausted,
	http.StatusInternalServerError:     codes.Internal,
	http.StatusNotImplemented:          codes.Unimplemented,
	http.StatusBadGateway:              codes.Unknown,
	http.StatusServiceUnavailable:      codes.Unavailable,
	http.StatusGatewayTimeout:          codes.DeadlineExceeded,
	http.StatusHTTPVersionNotSupported: codes.Unavailable,
	http.StatusRequestEntityTooLarge:   codes.OutOfRange,
	http.StatusUnsupportedMediaType:    codes.Canceled,
	http.StatusUnprocessableEntity:     codes.InvalidArgument,
	http.StatusMethodNotAllowed:        codes.Canceled,
	http.StatusPreconditionFailed:      codes.FailedPrecondition,
}
