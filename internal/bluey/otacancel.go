package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// OTACancel stops an OTA download.
func (b *Bluey) OTACancel(ctx context.Context, req *ep_bluetooth.OTACancelReq) (*ep_bluetooth.OTACancelResp, error) {
	if err := b.isNotAuthorized(); err != nil {
		return &ep_bluetooth.OTACancelResp{}, err
	}

	_, err := b.conn.OTACancel()
	if err != nil {
		return &ep_bluetooth.OTACancelResp{}, status.Error(codes.Unknown, err.Error())
	}

	return &ep_bluetooth.OTACancelResp{}, nil
}
