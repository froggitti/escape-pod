package interceptor

import (
	"context"
	"reflect"
	"testing"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/DDLbots/escape-pod/internal/license/interceptor/file"
	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
)

func TestInterceptor_Delete(t *testing.T) {
	lm, err := file.New(file.WithFilePath("./interceptor.delete.json"))
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
		req         *ep_license.DeleteReq
		setup       func(*Interceptor, []*format.License) error
		licenses    []*format.License
		want        *ep_license.DeleteResp
		wantErr     bool
	}{
		{
			name:        "should pass",
			interceptor: inter,
			licenses: []*format.License{
				{
					Email:   "testy@testerson.nil",
					Version: "1.0",
					Bot:     "vic:1234",
				},
			},
			setup:   addTestDevices,
			req:     &ep_license.DeleteReq{Bot: "vic:1234"},
			want:    &ep_license.DeleteResp{},
			wantErr: false,
		},
		{
			name:        "should fail; doesn't exist",
			interceptor: inter,
			licenses:    []*format.License{},
			setup: func(*Interceptor, []*format.License) error {
				return nil
			},
			req:     &ep_license.DeleteReq{Bot: "vic:bot_not_found"},
			want:    &ep_license.DeleteResp{},
			wantErr: true,
		},
		{
			name:        "should fail; empty request",
			interceptor: inter,
			licenses:    []*format.License{},
			setup: func(*Interceptor, []*format.License) error {
				return nil
			},
			want:    &ep_license.DeleteResp{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(tt.interceptor, tt.licenses); err != nil {
				t.Errorf("can't add test devices: %v", err)
			}

			got, err := tt.interceptor.Delete(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interceptor.Delete() = %v, want %v", got, tt.want)
			}

		})
	}
	t.Cleanup(func() {
		if err := lm.Drop(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}
