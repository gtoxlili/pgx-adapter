package pgxadapter

import (
	"context"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Adapter struct {
	store *store
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

func NewAdapter(ctx context.Context, db *pgxpool.Pool, opts ...Option) (*Adapter, error) {
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	adapter := &Adapter{
		store: newStore(db),
	}
	for _, opt := range opts {
		opt(adapter)
	}

	if err := adapter.store.initTable(ctx); err != nil {
		return nil, err
	}

	return adapter, nil
}
