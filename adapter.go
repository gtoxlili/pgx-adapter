package pgxadapter

import (
	"context"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Adapter struct {
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

func NewAdapter(ctx context.Context, db *pgxpool.Pool) (*Adapter, error) {
	panic("implement me")
}
