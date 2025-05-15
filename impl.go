package pgxadapter

import (
	"context"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

func (a *Adapter) LoadPolicyCtx(ctx context.Context, model model.Model) error {
	lines, err := a.store.selectAll(ctx)
	if err != nil {
		return err
	}
	for _, line := range lines {
		if err := persist.LoadPolicyArray(line, model); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) SavePolicyCtx(ctx context.Context, model model.Model) error {
	var args [][]string

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			args = append(args, genRule(ptype, rule))
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			args = append(args, genRule(ptype, rule))
		}
	}

	return a.store.deleteAndInsertAll(ctx, args)
}

func (a *Adapter) AddPolicyCtx(ctx context.Context, sec string, ptype string, rule []string) error {
	return a.store.insertRow(ctx, ptype, rule...)
}

func (a *Adapter) RemovePolicyCtx(ctx context.Context, sec string, ptype string, rule []string) error {
	return a.store.deleteRow(ctx, ptype, rule...)
}

func (a *Adapter) RemoveFilteredPolicyCtx(ctx context.Context, sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return a.store.deleteByPType(ctx, ptype)
}
