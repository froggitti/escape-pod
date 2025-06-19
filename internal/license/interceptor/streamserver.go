package interceptor

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StreamServerInterceptor checks the stream request for bot permission.
func (s *Interceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {

	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		s.RLock()
		if !s.authorize(stream.Context(), info.FullMethod) {
			return status.Errorf(codes.Unauthenticated, "this bot is not licensed")
		}
		s.RUnlock()
		return handler(srv, stream)
	}
}
