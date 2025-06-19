package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Status returns the BLE status.
func (b *Bluey) Status(ctx context.Context, req *ep_bluetooth.StatusReq) (*ep_bluetooth.StatusResp, error) {
	if err := b.isNotAuthorized(); err != nil {
		return &ep_bluetooth.StatusResp{}, err
	}

	sr, err := b.conn.GetStatus()
	if err != nil {
		return &ep_bluetooth.StatusResp{}, status.Error(codes.Unknown, err.Error())
	}

	r := ep_bluetooth.StatusResp{
		Status:   b.status,
		WifiSsid: sr.WifiSSID,
		Version:  sr.Version,
		Esn:      sr.ESN,
	}
	return &r, nil
}
