package pgxadapter

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/samber/lo"
)

func (a *Adapter) LoadFilteredPolicyCtx(ctx context.Context, model model.Model, filter interface{}) error {
	if lo.IsNil(filter) {
		return a.LoadPolicyCtx(ctx, model)
	}
	a.filter.Store(true)
	ft, ok := filter.(map[string][]string)
	if !ok {
		return fmt.Errorf("filter must be of type map[string][]string, got %T. Expected format: map[role][]fieldValues, where role is the key and fieldValues is a slice of strings. For unused attributes, use an empty string", filter)
	}
	var lines [][]string
	for k, v := range ft {
		tmp, err := a.store.selectWhere(ctx, k, 0, v...)
		if err != nil {
			return err
		}
		lines = append(lines, tmp...)
	}
	if len(lines) == 0 {
		return nil
	}
	for _, line := range lines {
		if err := persist.LoadPolicyArray(line, model); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) IsFilteredCtx(ctx context.Context) bool {
	return a.filter.Load()
}
