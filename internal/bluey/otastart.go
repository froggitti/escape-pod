package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	"github.com/digital-dream-labs/vector-bluetooth/ble"
)

// OTAStart tells the robot to fetch an OTA file.
func (b *Bluey) OTAStart(ctx context.Context, req *ep_bluetooth.OTAStartReq) (*ep_bluetooth.OTAStartResp, error) {
	if err := b.isNotAuthorized(); err != nil {
		return &ep_bluetooth.OTAStartResp{}, err
	}

	go b.runOTA(req)

	resp := ep_bluetooth.OTAStartResp{}

	return &resp, nil
}

func (b *Bluey) runOTA(req *ep_bluetooth.OTAStartReq) {
	resp, err := b.conn.OTAStart(req.Url)
	if err != nil {
		b.ch <- ble.StatusChannel{
			OTAStatus: &ble.StatusCounter{
				Error: err.Error(),
			},
		}
		return
	}

	b.ch <- ble.StatusChannel{
		OTAStatus: &ble.StatusCounter{
			Error: translateStatus(resp.Status),
		},
	}
}

func translateStatus(id int) string {
	statusmap := map[int]string{
		1:   "unknown status",
		2:   "ota in progress",
		3:   "ota completed",
		4:   "rebooting",
		5:   "ota error",
		10:  "unknown system error",
		200: "unexpected tar contents",
		201: "unhandled manifest version or feature",
		202: "boot control hal failure",
		203: "could not open url",
		204: "invalid file format",
		205: "decompress error",
		206: "block error",
		207: "imgdiff error",
		208: "i/o error",
		209: "signature validation error",
		210: "decryption error",
		211: "wrong base version",
		212: "subprocess exception",
		213: "wrong serial number",
		214: "dev/prod mismatch",
		215: "socket timeout error",
		216: "downgrade not allowed",
	}

	r, ok := statusmap[id]
	if !ok {
		return "other exception"
	}
	return r
}
