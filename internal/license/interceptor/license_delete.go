package interceptor

import (
	"context"
	"errors"

	"github.com/DDLbots/escape-pod/internal/license/interceptor/file"
	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Delete Deletes a license
func (s *Interceptor) Delete(ctx context.Context, req *ep_license.DeleteReq) (*ep_license.DeleteResp, error) {
	if req == nil || req.Bot == "" {
		return &ep_license.DeleteResp{}, status.Errorf(codes.InvalidArgument, "")
	}

	err := s.licenseManager.DeleteLicense(ctx, req.Bot)
	if err != nil {
		// log.WithFields(log.Fields{
		// 	"action": "delete",
		// 	"status": "failed",
		// }).Info(err)

		if errors.Is(err, file.ErrNotFound) {
			return &ep_license.DeleteResp{}, status.Errorf(codes.NotFound, "")
		}
		return &ep_license.DeleteResp{}, err
	}

	if err := s.load(); err != nil {
		// log.WithFields(log.Fields{
		// 	"action": "reload",
		// 	"status": "failed",
		// }).Info(err)
		return &ep_license.DeleteResp{}, status.Errorf(codes.Internal, "")
	}

	return &ep_license.DeleteResp{}, nil
}
