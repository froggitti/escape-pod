package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Scan returns a list of bluetooth devices
func (b *Bluey) Scan(ctx context.Context, req *ep_bluetooth.ScanReq) (*ep_bluetooth.ScanResp, error) {
	if err := b.isConnected(); err != nil {
		return &ep_bluetooth.ScanResp{}, err
	}

	r, err := b.conn.Scan()
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	ds := []*ep_bluetooth.Device{}
	for _, v := range r.Devices {
		d := ep_bluetooth.Device{
			Id:      int32(v.ID),
			Name:    v.Name,
			Address: v.Address,
		}
		ds = append(ds, &d)
	}

	resp := ep_bluetooth.ScanResp{
		Devices: ds,
	}

	return &resp, nil
}
