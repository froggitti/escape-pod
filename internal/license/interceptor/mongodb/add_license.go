package mongodb

import (
	"context"

	"github.com/DDLbots/escape-pod/internal/license/format"
)

func (lm *LicensesManager) AddLicense(ctx context.Context, in *format.Payload) error {
	lm.Lock()
	defer lm.Unlock()

	// cache, err := lm.load()
	// if err != nil {
	// 	return fmt.Errorf("load: %v", err)
	// }
	// if err := lm.close(); err != nil {
	// 	return err
	// }

	// for _, payload := range cache {
	// 	if cmp.Equal(payload, in) {
	// 		return ErrDocumentExists
	// 	}
	// }

	// cache = append(cache, in)

	// if err := lm.create(); err != nil {
	// 	return err
	// }

	// if err := lm.save(cache); err != nil {
	// 	return err
	// }

	// return lm.close()
	return nil
}
