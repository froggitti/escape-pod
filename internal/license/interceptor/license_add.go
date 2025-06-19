package interceptor

import (
	"context"
	"errors"

	"github.com/DDLbots/escape-pod/internal/license/interceptor/file"
	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Add adds a license
func (s *Interceptor) Add(ctx context.Context, req *ep_license.AddReq) (*ep_license.AddResp, error) {
	if req == nil {
		return &ep_license.AddResp{
			Response: "invalid license",
		}, status.Errorf(codes.InvalidArgument, "")
	}

	payload, err := s.validator.ValidateString(req.License)
	if err != nil {
		return &ep_license.AddResp{
			Response: "invalid license",
		}, status.Errorf(codes.InvalidArgument, "")
	}

	if err := s.licenseManager.AddLicense(ctx, payload); err != nil {
		if errors.Is(err, file.ErrDocumentExists) {
			return &ep_license.AddResp{
				Response: "already exists",
			}, status.Errorf(codes.AlreadyExists, "")
		}
		return nil, err
	}

	if err := s.load(); err != nil {
		// log.WithFields(log.Fields{
		// 	"action": "reload",
		// 	"status": "failed",
		// }).Info(err)
		return &ep_license.AddResp{}, status.Errorf(codes.Internal, "")
	}

	return &ep_license.AddResp{
		Response: "ok",
	}, nil

}
