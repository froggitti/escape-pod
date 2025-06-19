package bluey

import (
	"fmt"

	ep_bluetooth "github.com/DDLbots/internal-api/go/ep_bluetoothpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (b *Bluey) isNotConnectedOrAuthorized() error {
	switch b.status {
	case ep_bluetooth.Status_CONNECTED:
		return nil
	case ep_bluetooth.Status_AUTHORIZED:
		return nil
	default:
		return status.Error(codes.InvalidArgument, "this command cannot be performed on an already-authorized device")
	}
}

func (b *Bluey) isConnected() error {
	switch b.status {
	case ep_bluetooth.Status_CONNECTED:
		return nil
	default:
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("you may not run this command with connection of type %s", b.status.String()),
		)
	}
}

func (b *Bluey) isNotConnected() error {
	switch b.status {
	case ep_bluetooth.Status_NOT_CONNECTED:
		return nil
	case ep_bluetooth.Status_UNKNOWN:
		return nil
	default:
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("you may not run this command with connection of type %s", b.status.String()),
		)
	}
}

func (b *Bluey) isNotAuthorized() error {
	switch b.status {
	case ep_bluetooth.Status_AUTHORIZED:
		return nil
	default:
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("you may not run this command with connection of type %s", b.status.String()),
		)
	}
}

func (b *Bluey) isNotAuthenticated() error {
	switch b.status {
	case ep_bluetooth.Status_AUTHENTICATED:
		return nil
	default:
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("you may not run this command with connection of type %s", b.status.String()),
		)
	}
}

func (b *Bluey) hasConnection() error {
	switch b.status {
	case ep_bluetooth.Status_UNKNOWN:
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("you may not run this command with connection of type %s", b.status.String()),
		)
	case ep_bluetooth.Status_NOT_CONNECTED:
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("you may not run this command with connection of type %s", b.status.String()),
		)
	case ep_bluetooth.Status_BUSY:
		return status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("you may not run this command with connection of type %s", b.status.String()),
		)
	default:
		return nil
	}
}
