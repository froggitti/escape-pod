package bluey

import (
	"context"
	"io/ioutil"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
)

// ListLogs lists available logs!
func (b *Bluey) ListLogs(ctx context.Context, req *ep_bluetooth.ListLogsReq) (*ep_bluetooth.ListLogsResp, error) {
	list, err := ioutil.ReadDir(b.filepath)
	if err != nil {
		return &ep_bluetooth.ListLogsResp{}, err
	}

	files := []string{}
	for _, t := range list {
		files = append(files, t.Name())
	}

	resp := ep_bluetooth.ListLogsResp{
		Names: files,
	}
	return &resp, nil
}
