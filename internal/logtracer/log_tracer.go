package logtracer

import (
	ep_logtracer "github.com/DDLbots/internal-api/go/ep_logtracerpb"
	grpc "google.golang.org/grpc"
)

// Deprecated
type LogTracer struct {
	ep_logtracer.UnimplementedLogTracerServer
}

// Start starts the logtrace service
func Start(transport *grpc.Server) {
	s := &LogTracer{}

	ep_logtracer.RegisterLogTracerServer(transport, s)
}
