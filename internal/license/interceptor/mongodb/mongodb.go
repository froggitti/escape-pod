package mongodb

import (
	"context"
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrDocumentExists = errors.New("document exists")
	ErrNotFound       = errors.New("not found")
)

type LicensesManager struct {
	sync.RWMutex
	collectionName string
	db             *mongo.Database
	collection     *mongo.Collection
}

func New(opts ...Option) (*LicensesManager, error) {
	lm := &LicensesManager{}
	for _, opt := range opts {
		if err := opt.applyOption(lm); err != nil {
			return nil, err
		}
	}
	return lm, lm.IsValid()
}
func (lm *LicensesManager) IsValid() error {
	if lm.collectionName == "" {
		return errors.New("empty collection name")
	}

	if lm.db == nil {
		return errors.New("empty database")
	}

	lm.collection = lm.db.Collection(lm.collectionName)

	return nil
}

func (lm *LicensesManager) Drop(context.Context) error {
	lm.Lock()
	defer lm.Unlock()
	// return lm.drop()
	return nil
}

// func (lm *LicensesManager) drop() error {
// 	return os.Remove(lm.filePath)
// }

func (lm *LicensesManager) Purge(context.Context) error {
	lm.Lock()
	defer lm.Unlock()
	// if err := lm.drop(); err != nil {
	// 	return err
	// }
	// return lm.create()
	return nil
}

// func (lm *LicensesManager) create() error {
// 	var err error
// 	lm.file, err = os.Create(lm.filePath)
// 	return err
// }

// func (lm *LicensesManager) load() ([]*format.Payload, error) {
// 	var err error
// 	lm.file, err = os.OpenFile(lm.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
// 	if err != nil {
// 		return nil, fmt.Errorf("open file: %v", err)
// 	}
// 	cache := []*format.Payload{}
// 	wr := json.NewDecoder(lm.file)

// 	if err := wr.Decode(&cache); err != nil {
// 		if !errors.Is(err, io.EOF) {
// 			return nil, fmt.Errorf("decode: %v", err)
// 		}
// 	}

// 	return cache, nil
// }

// func (lm *LicensesManager) close() error {
// 	if lm.file != nil {
// 		if err := lm.file.Close(); err != nil {
// 			return err
// 		}
// 		lm.file = nil
// 	}
// 	return nil
// }

// func (lm *LicensesManager) save(cache []*format.Payload) error {
// 	if lm.file == nil {
// 		return errors.New("file not found")
// 	}
// 	wr := json.NewEncoder(lm.file)

// 	if err := wr.Encode(cache); err != nil {
// 		return err
// 	}

// 	return nil
// }
