package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Close turns off the BLE connection.
func (b *Bluey) Close(ctx context.Context, req *ep_bluetooth.CloseReq) (*ep_bluetooth.CloseResp, error) {
	if b.conn == nil {
		return &ep_bluetooth.CloseResp{}, status.Error(codes.InvalidArgument, "no active connection")
	}

	if err := b.conn.Close(); err != nil {
		return &ep_bluetooth.CloseResp{}, status.Error(codes.Internal, "cannot close connection: "+err.Error())
	}

	b.status = ep_bluetooth.Status_NOT_CONNECTED
	b.conn = nil

	return &ep_bluetooth.CloseResp{}, status.Error(codes.OK, "ble connection closed")
}
