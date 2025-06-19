package mongodb

import (
	"context"
	"testing"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLicensesManager_ListLicenses(t *testing.T) {
	lm, err := New(
	// WithFilePath("./list.licenses.json"),
	)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		lm      *LicensesManager
		args    args
		want    []*format.Payload
		wantErr bool
	}{
		{
			name: "should pass",
			lm:   lm,
			args: args{ctx: context.Background()},
			want: []*format.Payload{
				{License: &format.License{Bot: "test add bod"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, license := range tt.want {
				if err := tt.lm.AddLicense(context.Background(), license); err != nil {
					t.Fatal(err)
				}
			}

			got, err := tt.lm.ListLicenses(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("LicensesManager.ListLicenses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opts := cmp.Options{
				cmpopts.IgnoreUnexported(),
			}

			if !cmp.Equal(got, tt.want, opts) {
				t.Errorf("LicensesManager.ListLicenses() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Cleanup(func() {
		if err := lm.Drop(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}
