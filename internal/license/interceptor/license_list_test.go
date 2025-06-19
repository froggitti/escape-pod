package interceptor

import (
	"context"
	"testing"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/DDLbots/escape-pod/internal/license/interceptor/file"
	"github.com/DDLbots/escape-pod/internal/license/issuer"
	ep_license "github.com/DDLbots/internal-api/go/ep_licensepb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestInterceptor_List(t *testing.T) {
	lm, err := file.New(file.WithFilePath("./interceptor.list.json"))
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
		licenses    []*format.License
		setup       func(*Interceptor, []*format.License) error
		want        *ep_license.ListResp
		wantErr     bool
	}{
		{
			name:        "should pass; with zero robots",
			interceptor: inter,
			licenses: []*format.License{
				{
					Email:   "testy@testerson.nil",
					Version: "1.0",
					Bot:     "vic:1234",
				},
			},
			setup: func(*Interceptor, []*format.License) error {
				return nil
			},
			want: &ep_license.ListResp{
				Bots: []string{},
			},
			wantErr: false,
		},
		{
			name:        "should pass; one",
			interceptor: inter,
			licenses: []*format.License{
				{
					Email:   "testy@testerson.nil",
					Version: "1.0",
					Bot:     "vic:1234",
				},
			},
			setup:   addTestDevices,
			want:    &ep_license.ListResp{Bots: []string{"vic:1234"}},
			wantErr: false,
		},
		{
			name:        "should pass; two",
			interceptor: inter,
			licenses: []*format.License{
				{
					Email:   "testy@testerson.nil",
					Version: "1.0",
					Bot:     "vic:test_two_first",
				},
				{
					Email:   "testy@testerson.nil",
					Version: "1.0",
					Bot:     "vic:test_two_second",
				},
			},
			setup:   addTestDevices,
			want:    &ep_license.ListResp{Bots: []string{"vic:test_two_first", "vic:test_two_second"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(tt.interceptor, tt.licenses); err != nil {
				t.Errorf("can't add test devices: %v", err)
			}

			got, err := tt.interceptor.List(context.Background(), &ep_license.ListReq{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Interceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmp.Options{
				cmpopts.IgnoreUnexported(ep_license.ListResp{}),
			}

			if !cmp.Equal(got, tt.want, opts...) {
				t.Errorf("%+v\n", cmp.Diff(got, tt.want, opts...))
				t.Errorf("Interceptor.List() = %v, want %v", got, tt.want)
			}
		})
		if err := lm.Purge(context.Background()); err != nil {
			t.Fatal(err)
		}

	}
	t.Cleanup(func() {
		if err := lm.Drop(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}

func addTestDevices(i *Interceptor, args []*format.License) error {
	iss := issuer.New()
	for v := range args {
		lic, err := iss.Generate(args[v])
		if err != nil {
			return err
		}

		if _, err := i.Add(
			context.Background(),
			&ep_license.AddReq{
				License: lic,
			},
		); err != nil {
			return err
		}
	}
	return nil
}
