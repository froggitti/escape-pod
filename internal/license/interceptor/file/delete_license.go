package file

import (
	"context"
	"fmt"

	"github.com/DDLbots/escape-pod/internal/license/format"
)

func remove(slice []*format.Payload, s int) []*format.Payload {
	return append(slice[:s], slice[s+1:]...)
}

func (lm *LicensesManager) DeleteLicense(ctx context.Context, bot string) error {
	lm.Lock()
	defer lm.Unlock()

	cache, err := lm.load()
	if err != nil {
		return err
	}

	if err := lm.close(); err != nil {
		return err
	}

	var oerr error
	found := false
	for i, payload := range cache {
		if payload.License.Bot == bot {
			found = true
			cache = remove(cache, i)
			break
		}
	}

	// TODO: see what this looks like
	if !found {
		oerr = fmt.Errorf("delete: %v", ErrNotFound)
	}

	if err := lm.create(); err != nil {
		return err
	}

	if err := lm.save(cache); err != nil {
		oerr = fmt.Errorf("save: %v", err)
	}

	if err := lm.close(); err != nil {
		oerr = fmt.Errorf("save: %v", err)
	}
	lm.Debug("end of delete")
	return oerr
}
