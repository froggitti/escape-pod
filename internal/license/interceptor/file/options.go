package file

import "github.com/DDLbots/go-logger"

type Option interface {
	applyOption(*LicensesManager) error
}
type applyOptionFunc func(*LicensesManager) error

func (f applyOptionFunc) applyOption(lm *LicensesManager) error {
	return f(lm)
}

func WithFilePath(filePath string) Option {
	return applyOptionFunc(func(lm *LicensesManager) error {
		lm.filePath = filePath
		return nil
	})
}

func WithDebugger(debugger logger.Debugger) Option {
	return applyOptionFunc(func(lm *LicensesManager) error {
		lm.debugger = debugger
		return nil
	})
}
