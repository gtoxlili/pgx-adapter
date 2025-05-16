package pgxadapter

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/samber/lo"
	"strings"
)

func (a *Adapter) LoadFilteredPolicyCtx(ctx context.Context, model model.Model, filter interface{}) error {
	if lo.IsNil(filter) {
		return a.LoadPolicyCtx(ctx, model)
	}
	a.filter.Store(true)
	ft, ok := filter.(map[string][][]string)
	if !ok {
		return fmt.Errorf("filter must be of type map[string][][]string, got %T. "+
			"Expected format: map[role][][]fieldValues, "+
			"where role is the key (like 'p', 'g', 'g2') and fieldValues is a 2D slice of strings representing multiple conditions with OR relationship between them. "+
			"For unused attributes in conditions, use an empty string", filter)
	}
	var lines [][]string
	for k, v := range ft {
		for _, vv := range v {
			tmp, err := a.store.selectWhere(ctx, k, 0, vv...)
			if err != nil {
				return err
			}
			lines = append(lines, tmp...)
		}
	}
	lines = lo.UniqBy(lines, func(line []string) string {
		return strings.Join(line, ",")
	})
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
