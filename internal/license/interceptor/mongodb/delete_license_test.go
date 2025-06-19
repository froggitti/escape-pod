package mongodb

import (
	"context"
	"testing"

	"github.com/DDLbots/escape-pod/internal/license/format"
)

func TestLicensesManager_DeleteLicense(t *testing.T) {
	lm, err := New(
	// WithCollectionName("delete"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := lm.AddLicense(context.Background(), &format.Payload{License: &format.License{Bot: "delete_test_bot"}}); err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx context.Context
		bot string
	}
	tests := []struct {
		name    string
		lm      *LicensesManager
		args    args
		wantErr bool
	}{
		{
			name: "should pass",
			lm:   lm,
			args: args{
				ctx: context.Background(),
				bot: "delete_test_bot",
			},
		},
		{
			name: "should fail; bot not found",
			lm:   lm,
			args: args{
				ctx: context.Background(),
				bot: "delete_test_bot_not_found",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.lm.DeleteLicense(tt.args.ctx, tt.args.bot); (err != nil) != tt.wantErr {
				t.Errorf("LicensesManager.DeleteLicense() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}
	t.Cleanup(func() {
		if err := lm.Drop(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}
