package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/DDLbots/escape-pod/internal/license/format"
	"github.com/DDLbots/go-logger"
)

var (
	ErrDocumentExists = errors.New("document exists")
	ErrNotFound       = errors.New("not found")
)

type LicensesManager struct {
	sync.RWMutex
	filePath string
	file     io.ReadWriteCloser
	debugger logger.Debugger
}

func New(opts ...Option) (*LicensesManager, error) {
	lm := &LicensesManager{
		filePath: filepath.Join(os.TempDir(), "licenses"),
	}
	for _, opt := range opts {
		if err := opt.applyOption(lm); err != nil {
			return nil, err
		}
	}

	lm.load()
	return lm, nil
}

func (lm *LicensesManager) Drop(context.Context) error {
	lm.Lock()
	defer lm.Unlock()
	return lm.drop()
}

func (lm *LicensesManager) drop() error {
	lm.Debug("drop")

	return os.Remove(lm.filePath)
}

func (lm *LicensesManager) Purge(context.Context) error {
	lm.Lock()
	defer lm.Unlock()
	// if err := lm.drop(); err != nil {
	// 	return err
	// }
	return lm.create()
}

func (lm *LicensesManager) create() error {
	lm.Debug("create")

	var err error
	lm.file, err = os.Create(lm.filePath)
	return err
}

func (lm *LicensesManager) load() ([]*format.Payload, error) {
	var err error
	lm.Debug("start open file")

	lm.file, err = os.OpenFile(lm.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}

	cache := []*format.Payload{}
	wr := json.NewDecoder(lm.file)

	lm.Debug("start decode")
	if err := wr.Decode(&cache); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("decode: %v", err)
		}
	}

	return cache, nil
}

func (lm *LicensesManager) close() error {
	lm.Debug("close")
	if lm.file != nil {
		if err := lm.file.Close(); err != nil {
			return err
		}
		lm.file = nil
	}
	return nil
}

func (lm *LicensesManager) save(cache []*format.Payload) error {
	lm.Debug("save")

	if lm.file == nil {
		return errors.New("file not found")
	}
	wr := json.NewEncoder(lm.file)

	lm.Debug("start encode")

	if err := wr.Encode(cache); err != nil {
		return err
	}

	return nil
}

func (lm *LicensesManager) Debug(msg string, args ...any) {
	if lm.debugger != nil {
		lm.debugger.Debug(msg, args...)
	}
}
func (lm *LicensesManager) Debugf(msg string, args ...any) {
	if lm.debugger != nil {
		lm.debugger.Debug(fmt.Sprintf(msg, args...))
	}
}
