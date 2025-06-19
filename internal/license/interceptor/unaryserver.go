package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor checks the unary request for bot permission.
func (s *Interceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		s.RLock()

		if !s.authorize(ctx, info.FullMethod) {
			return nil, status.Errorf(codes.Unauthenticated, "this bot is not licensed")
		}
		s.RUnlock()

		return handler(ctx, req)
	}
}
