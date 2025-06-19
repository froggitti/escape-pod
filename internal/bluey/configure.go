package bluey

import (
	"context"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	"github.com/digital-dream-labs/vector-bluetooth/ble"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Configure sends a configuration message to the robot.
func (b *Bluey) Configure(ctx context.Context, req *ep_bluetooth.ConfigureReq) (*ep_bluetooth.ConfigureResp, error) {
	if err := b.isNotAuthenticated(); err != nil {
		return &ep_bluetooth.ConfigureResp{}, err
	}

	if err := b.conn.ConfigureSettings(
		&ble.VectorSettings{
			Timezone:           req.Timezone,
			DefaultLocation:    req.DefaultLocation,
			Locale:             req.Locale,
			AllowDataAnalytics: req.AllowDataAnalytics,
			MetricDistance:     req.MetricDistance,
			MetricTemperature:  req.MetricTemperatire,
			AlexaOptIn:         req.AlexaOptIn,
			ButtonWakeword:     req.ButtonWakeword,
			Clock24Hour:        req.Clock_24Hour,
		},
	); err != nil {
		return &ep_bluetooth.ConfigureResp{}, status.Error(codes.Unknown, err.Error())
	}

	return &ep_bluetooth.ConfigureResp{}, nil
}
