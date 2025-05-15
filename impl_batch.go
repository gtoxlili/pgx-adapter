package pgxadapter

import "context"

func (a *Adapter) AddPoliciesCtx(ctx context.Context, sec string, ptype string, rules [][]string) error {
	return a.store.batchInsert(ctx, ptype, rules)
}

func (a *Adapter) RemovePoliciesCtx(ctx context.Context, sec string, ptype string, rules [][]string) error {
	return a.store.batchDelete(ctx, ptype, rules)
}
