package pgxadapter

import (
	"context"
	"database/sql/driver"
	"github.com/casbin/casbin/v2/persist"
	"go.uber.org/atomic"
)

type Adapter struct {
	store  *store
	filter *atomic.Bool
}

// the supported for Casbin interfaces.
var (
	_ persist.Adapter                 = new(Adapter)
	_ persist.ContextAdapter          = new(Adapter)
	_ persist.FilteredAdapter         = new(Adapter)
	_ persist.ContextFilteredAdapter  = new(Adapter)
	_ persist.BatchAdapter            = new(Adapter)
	_ persist.ContextBatchAdapter     = new(Adapter)
	_ persist.UpdatableAdapter        = new(Adapter)
	_ persist.ContextUpdatableAdapter = new(Adapter)
)

type Option func(*Adapter)

func WithFieldCount(fieldCount int) Option {
	return func(a *Adapter) {
		a.store.setFieldCount(fieldCount)
	}
}

func WithTableName(tableName string) Option {
	return func(a *Adapter) {
		a.store.setTableName(tableName)
	}
}

func NewAdapter(ctx context.Context, db interface {
	driver.Pinger
	Storer
}, opts ...Option) (*Adapter, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	adapter := &Adapter{
		store:  newStore(db),
		filter: atomic.NewBool(false),
	}
	for _, opt := range opts {
		opt(adapter)
	}

	if err := adapter.store.initTable(ctx); err != nil {
		return nil, err
	}

	return adapter, nil
}
