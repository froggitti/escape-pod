package interceptor

import (
	"context"

	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// List Lists a license
func (s *Interceptor) List(ctx context.Context, req *ep_license.ListReq) (*ep_license.ListResp, error) {
	s.RLock()
	defer s.RUnlock()

	bots, err := s.licenseManager.ListBots(ctx)
	if err != nil {
		return &ep_license.ListResp{}, status.Error(codes.Internal, "INTERNAL_ERROR")
	}

	return &ep_license.ListResp{
		Bots: bots,
	}, nil
}
