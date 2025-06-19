package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// WifiConnect tells the vector to connect to a specified wifi network.
func (b *Bluey) WifiConnect(ctx context.Context, req *ep_bluetooth.WifiConnectReq) (*ep_bluetooth.WifiConnectResp, error) {
	if err := b.isNotAuthorized(); err != nil {
		return &ep_bluetooth.WifiConnectResp{}, err
	}

	r, err := b.conn.WifiConnect(req.WifiSsid, req.WifiPassword, 5, int(req.Authtype))
	if err != nil {
		return &ep_bluetooth.WifiConnectResp{}, status.Error(codes.Unknown, err.Error())
	}

	resp := ep_bluetooth.WifiConnectResp{
		State:  int32(r.State),
		Result: int32(r.Result),
	}

	return &resp, nil
}
