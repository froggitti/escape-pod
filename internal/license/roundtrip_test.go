// +build !production

package license

import (
	"reflect"
	"testing"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/DDLbots/escape-pod/internal/license/issuer"
	"github.com/DDLbots/escape-pod/internal/license/validator"
)

func TestIssuer_Generate(t *testing.T) {
	tests := []struct {
		name    string
		args    *format.License
		check   bool
		wantErr bool
	}{
		{
			name: "passing",
			args: &format.License{
				Email:   "test@test.com",
				Version: "1.0",
				Bot:     "vic:00a10c90",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iss := issuer.New()
			val := validator.New()

			got, err := iss.Generate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Licensor.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			v, err := val.ValidateString(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Licensor.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.args, v.License) {
				t.Errorf("Validator.Validate() = %v, want %v", got, v)
			}
		})
	}
}
