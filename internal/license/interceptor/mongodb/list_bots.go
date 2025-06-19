package mongodb

import (
	"context"
)

func (lm *LicensesManager) ListBots(ctx context.Context) ([]string, error) {
	lm.Lock()
	defer lm.Unlock()

	// cache, err := lm.load()
	// if err != nil {
	// 	return nil, err
	// }

	// out := []string{}

	// for _, payload := range cache {
	// 	if payload.License != nil && payload.License.Bot != "" {
	// 		out = append(out, payload.License.Bot)
	// 	}
	// }

	// return out, lm.close()
	return nil, nil
}
