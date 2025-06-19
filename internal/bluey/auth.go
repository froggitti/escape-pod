package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Auth is used to send the cloud token to the robot.
func (b *Bluey) Auth(ctx context.Context, req *ep_bluetooth.AuthReq) (*ep_bluetooth.AuthResp, error) {
	if err := b.isNotAuthorized(); err != nil {
		return &ep_bluetooth.AuthResp{}, err
	}

	r, err := b.conn.Auth(req.Token)
	if err != nil {
		return &ep_bluetooth.AuthResp{}, status.Error(codes.Unknown, err.Error())
	}

	if !r.Success {
		return &ep_bluetooth.AuthResp{}, status.Error(codes.Unauthenticated, "authentication failed")
	}

	resp := ep_bluetooth.AuthResp{
		//Status: int32(r.Status),
		Success: r.Success,
	}

	b.status = ep_bluetooth.Status_AUTHENTICATED

	return &resp, nil
}
