package interceptor

import (
	"context"
	"errors"
)

func (s *Interceptor) load() error {
	s.Lock()
	defer s.Unlock()

	licenses, err := s.licenseManager.ListLicenses(context.Background())
	if err != nil {
		return err
	}

	if len(licenses) == 0 {
		// log.WithFields(log.Fields{
		// 	"license error": "zero bots are authorized -- you'll need to add your license key",
		// }).Warn("err")
	}

	bots := map[string]struct{}{}

	for _, license := range licenses {
		if err := s.validator.ValidatePayload(license); err != nil {
			return errors.New("invalid license key detected, disabling ALL bots -- that'll learn you")
		}
		bots[license.License.Bot] = struct{}{}
	}

	s.bots = bots

	return nil
}
