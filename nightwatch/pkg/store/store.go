package store

import (
	"context"
	"sync"

	"gorm.io/gorm"
)

var (
	once  sync.Once
	Store *datastore
)

type txKey struct{}

type IStore interface {
	TX(context.Context, func(ctx context.Context) error) error
	Tasks() TaskStore
}

type datastore struct {
	core *gorm.DB
}

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) *datastore {
	once.Do(func() {
		Store = &datastore{db}
	})

	return Store
}

func (ds *datastore) Core(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok {
		return tx
	}

	return ds.core
}

func (ds *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return ds.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, txKey{}, tx)
			return fn(ctx)
		},
	)
}

func (ds *datastore) Tasks() TaskStore {
	return newTaskStore(ds)
}
