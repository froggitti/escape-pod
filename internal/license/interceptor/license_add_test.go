package interceptor

import (
	"context"
	"log"
	"testing"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/DDLbots/escape-pod/internal/license/interceptor/file"
	"github.com/DDLbots/escape-pod/internal/license/issuer"
	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
)

func TestInterceptor_Add(t *testing.T) {

	lm, err := file.New(file.WithFilePath("./interceptor.add.json"))
	if err != nil {
		t.Fatal(err)
	}

	inter, err := newInterceptor(lm)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		interceptor *Interceptor
		setup       func(*format.License) (string, error)
		license     *format.License
		validator   func() (*format.Payload, error)
		wantErr     bool
	}{
		{
			name:        "should pass",
			interceptor: inter,
			setup: func(req *format.License) (string, error) {
				i := issuer.New()
				return i.Generate(req)
			},
			license: &format.License{
				Email:   "testy@testerson.nil",
				Version: "1.0",
				Bot:     "vic:1234",
			},
		},
		{
			name:        "invalid license",
			interceptor: inter,
			setup: func(req *format.License) (string, error) {
				return "12314q2512525154", nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			license, err := tt.setup(tt.license)
			if err != nil {
				log.Fatal(err)
			}

			_, err = tt.interceptor.Add(
				context.Background(),
				&ep_license.AddReq{
					License: license,
				},
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interceptor.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
	t.Cleanup(func() {
		if err := lm.Drop(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}
