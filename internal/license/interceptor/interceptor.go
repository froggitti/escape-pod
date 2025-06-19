package interceptor

import (
	"context"
	"sync"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/DDLbots/escape-pod/internal/license/validator"
	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
)

type LicenseManager interface {
	AddLicense(ctx context.Context, payload *format.Payload) error
	DeleteLicense(ctx context.Context, bot string) error
	ListLicenses(ctx context.Context) ([]*format.Payload, error)
	ListBots(ctx context.Context) ([]string, error)
	Drop(context.Context) error
	Purge(context.Context) error
}

type StringValidator interface {
	ValidateString(string) (*format.Payload, error)
	ValidatePayload(req *format.Payload) error
}

// Interceptor is the struct to which the functions belong.
type Interceptor struct {
	sync.RWMutex
	bots map[string]struct{}
	key  []byte

	validator StringValidator

	licenseManager LicenseManager
	ep_license.UnimplementedLicenseManagerServer
}

// New reads the license file and returns an interceptor.
func New(signingKey string, licenseManager LicenseManager) (*Interceptor, error) {
	i := &Interceptor{
		validator:      validator.New(),
		key:            []byte(signingKey),
		licenseManager: licenseManager,
	}

	if err := i.load(); err != nil {
		return nil, err
	}

	return i, nil
}

// func Start(signingKey string, licenseManager LicenseManager, transport *grpcserver.Server) (*Interceptor, error) {
// 	i, err := New(signingKey, licenseManager)
// 	if err != nil {
// 		return nil, err
// 	}

// 	registerService(i, transport)

// 	return i, nil
// }

// func registerService(in ep_license.LicenseManagerServer, transport *grpcserver.Server) {
// 	if err := transport.RegisterHTTPService(
// 		[]func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
// 			ep_license.RegisterLicenseManagerHandlerFromEndpoint,
// 		},
// 	); err != nil {
// 		log.Fatal(err)
// 	}

// 	ep_license.RegisterLicenseManagerServer(transport.Transport(), in)
// }
