package bluey

import (
	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"

	"github.com/digital-dream-labs/vector-bluetooth/ble"
)

// Bluey is the configuration struct
type Bluey struct {
	conn     *ble.VectorBLE
	ch       chan ble.StatusChannel
	status   ep_bluetooth.Status
	filepath string

	ep_bluetooth.UnimplementedBlueyServer
}

// New returns a populated struct
func New(dir string, ch chan ble.StatusChannel) (*Bluey, error) {
	return &Bluey{
		status:   ep_bluetooth.Status_NOT_CONNECTED,
		filepath: dir,
		ch:       ch,
	}, nil
}

// // Start starts the bluey management service
// func Start(transport *grpcserver.Server, dir string, ch chan ble.StatusChannel) {
// 	if err := transport.RegisterHTTPService(
// 		[]func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
// 			ep_bluetooth.RegisterBlueyHandlerFromEndpoint,
// 		},
// 	); err != nil {
// 		log.Fatal(err)
// 	}

// 	b, err := New(dir, ch)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ep_bluetooth.RegisterBlueyServer(transport.Transport(), b)
// }
