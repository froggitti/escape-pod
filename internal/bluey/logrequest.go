package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// FetchLogs backround-queues a log fetch operation
func (b *Bluey) FetchLogs(ctx context.Context, req *ep_bluetooth.FetchLogsReq) (*ep_bluetooth.FetchLogsResp, error) {
	if err := b.hasConnection(); err != nil {
		return &ep_bluetooth.FetchLogsResp{}, status.Error(codes.Unknown, err.Error())
	}

	go b.backgroundLogRequest()
	return &ep_bluetooth.FetchLogsResp{}, status.Error(codes.OK, "ble connection closed")
}

func (b *Bluey) backgroundLogRequest() {
	t := b.status
	b.status = ep_bluetooth.Status_BUSY
	_, _ = b.conn.DownloadLogs()
	b.status = t
}
