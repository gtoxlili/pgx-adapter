package pgxadapter

import "github.com/casbin/casbin/v2/persist"

type Adapter struct {
}

var _ persist.Adapter = (*Adapter)(nil)

var _ persist.FilteredAdapter = (*Adapter)(nil)

var _ persist.BatchAdapter = (*Adapter)(nil)

var _ persist.UpdatableAdapter = (*Adapter)(nil)

var _ persist.ContextAdapter = (*Adapter)(nil)

var _ persist.ContextBatchAdapter = (*Adapter)(nil)

var _ persist.ContextUpdatableAdapter = (*Adapter)(nil)

var _ persist.ContextFilteredAdapter = (*Adapter)(nil)
