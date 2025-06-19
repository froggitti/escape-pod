package mongodb

import (
	"context"
	"testing"

	"github.com/DDLbots/escape-pod/internal/license/format"
)

func TestLicensesManager_AddLicense(t *testing.T) {
	lm, err := New(
	// WithFilePath("./add.licenses.json"),
	)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		ctx context.Context
		in  *format.Payload
	}
	tests := []struct {
		name    string
		lm      *LicensesManager
		args    args
		wantErr bool
	}{
		{
			name: "should add",
			lm:   lm,
			args: args{
				ctx: context.Background(),
				in:  &format.Payload{License: &format.License{Email: "test_email", Bot: "test_bot"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.lm.AddLicense(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("LicensesManager.AddLicense() error = %v, wantErr %v", err, tt.wantErr)
			}

			bots, err := tt.lm.ListBots(tt.args.ctx)
			if err != nil {
				t.Fatal(err)
			}
			for _, bot := range bots {
				t.Logf("bot: %s", bot)
			}
		})
		t.Cleanup(func() {
			if err := lm.Drop(context.Background()); err != nil {
				t.Fatal(err)
			}
		})
	}
}
