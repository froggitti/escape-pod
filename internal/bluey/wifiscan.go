package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// WifiScan returns a list of available WIFI networks.
func (b *Bluey) WifiScan(ctx context.Context, req *ep_bluetooth.WifiScanReq) (*ep_bluetooth.WifiScanResp, error) {
	if err := b.isNotAuthorized(); err != nil {
		return &ep_bluetooth.WifiScanResp{}, err
	}

	nw, err := b.conn.WifiScan()
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	networks := []*ep_bluetooth.WifiNetwork{}

	for _, v := range nw.Networks {
		networks = append(
			networks,
			&ep_bluetooth.WifiNetwork{
				WifiSsid:       v.WifiSSID,
				SignalStrength: int32(v.SignalStrength),
				AuthType:       int32(v.AuthType),
				Hidden:         v.Hidden,
			},
		)
	}

	resp := ep_bluetooth.WifiScanResp{
		Networks: networks,
	}

	return &resp, nil
}
