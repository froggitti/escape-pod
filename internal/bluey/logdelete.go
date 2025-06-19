package bluey

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// DeleteLogs deletes logs!
func (b *Bluey) DeleteLogs(ctx context.Context, req *ep_bluetooth.DeleteLogsReq) (*ep_bluetooth.DeleteLogsResp, error) {
	if err := os.Remove(
		filepath.Clean(
			fmt.Sprintf(
				"%s/%s",
				b.filepath,
				req.Name,
			),
		),
	); err != nil {
		return &ep_bluetooth.DeleteLogsResp{}, status.Error(codes.Internal, "cannot delete file path: "+err.Error())
	}
	return &ep_bluetooth.DeleteLogsResp{}, nil
}
