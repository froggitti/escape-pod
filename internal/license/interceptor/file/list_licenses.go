package file

import (
	"context"

	"github.com/DDLbots/escape-pod/internal/license/format"
)

func (lm *LicensesManager) ListLicenses(ctx context.Context) ([]*format.Payload, error) {
	lm.RLock()
	defer lm.RUnlock()

	cache, err := lm.load()
	if err != nil {
		return nil, err
	}

	if err := lm.close(); err != nil {
		return nil, err
	}

	lm.Debug("end of list licenses")
	return cache, nil
}
