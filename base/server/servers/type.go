package servers

type ProtocolType string

const (
	GrpcServer ProtocolType = "grpc"
	HttpServer ProtocolType = "http"
)
