package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type Option interface {
	applyOption(*LicensesManager) error
}
type applyOptionFunc func(*LicensesManager) error

func (f applyOptionFunc) applyOption(lm *LicensesManager) error {
	return f(lm)
}

func WithCollectionName(collectionName string) Option {
	return applyOptionFunc(func(lm *LicensesManager) error {
		lm.collectionName = collectionName
		return nil
	})
}

func WithDB(db *mongo.Database) Option {
	return applyOptionFunc(func(lm *LicensesManager) error {
		lm.db = db
		return nil
	})
}
