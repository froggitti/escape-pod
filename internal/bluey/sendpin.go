package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// SendPin sends the pin that the user should retrieve from the screen.
func (b *Bluey) SendPin(ctx context.Context, req *ep_bluetooth.SendPinReq) (*ep_bluetooth.SendPinResp, error) {
	if err := b.isConnected(); err != nil {
		return &ep_bluetooth.SendPinResp{}, err
	}

	if err := b.conn.SendPin(req.Pin); err != nil {
		return &ep_bluetooth.SendPinResp{}, status.Error(codes.Unknown, "invalid pin")
	}

	_, err := b.conn.GetStatus()
	if err != nil {
		return &ep_bluetooth.SendPinResp{}, status.Error(codes.Unknown, err.Error())
	}

	b.status = ep_bluetooth.Status_AUTHORIZED

	return &ep_bluetooth.SendPinResp{}, nil
}
