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
	_ persist.Adapter                 = (*Adapter)(nil)
	_ persist.ContextAdapter          = (*Adapter)(nil)
	_ persist.FilteredAdapter         = (*Adapter)(nil)
	_ persist.ContextFilteredAdapter  = (*Adapter)(nil)
	_ persist.BatchAdapter            = (*Adapter)(nil)
	_ persist.ContextBatchAdapter     = (*Adapter)(nil)
	_ persist.UpdatableAdapter        = (*Adapter)(nil)
	_ persist.ContextUpdatableAdapter = (*Adapter)(nil)
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

// WithNoRowsAffectedError 当 受影响的行数 为 0 时，返回的错误（默认为 nil）
func WithNoRowsAffectedError(err error) Option {
	return func(a *Adapter) {
		a.store.setNoRowsAffectedError(err)
	}
}

func NewAdapter(ctx context.Context, db interface {
	driver.Pinger
	Commander
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
