package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Connect connects to a vector device.
func (b *Bluey) Connect(ctx context.Context, req *ep_bluetooth.ConnectReq) (*ep_bluetooth.ConnectResp, error) {
	if err := b.isConnected(); err != nil {
		return &ep_bluetooth.ConnectResp{}, err
	}

	if err := b.conn.Connect(int(req.Id)); err != nil {
		return &ep_bluetooth.ConnectResp{}, status.Error(codes.Unknown, err.Error())
	}

	b.status = ep_bluetooth.Status_CONNECTED

	return &ep_bluetooth.ConnectResp{
		Status: b.status,
	}, nil
}
