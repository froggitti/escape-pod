package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	"github.com/digital-dream-labs/vector-bluetooth/ble"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Init turns on the BLE connection.
func (b *Bluey) Init(ctx context.Context, req *ep_bluetooth.InitReq) (*ep_bluetooth.InitResp, error) {
	if err := b.isNotConnected(); err != nil {
		return &ep_bluetooth.InitResp{}, err
	}

	conn, err := ble.New(
		ble.WithLogDirectory(b.filepath),
		ble.WithStatusChan(b.ch),
	)
	if err != nil {
		return &ep_bluetooth.InitResp{}, status.Error(codes.Unknown, err.Error())
	}

	b.conn = conn
	b.status = ep_bluetooth.Status_CONNECTED

	return &ep_bluetooth.InitResp{
		Status: b.status,
	}, status.Error(codes.OK, "")
}
