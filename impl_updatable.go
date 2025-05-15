package pgxadapter

import (
	"context"
	"errors"
)

func (a *Adapter) UpdatePolicyCtx(ctx context.Context, sec string, ptype string, oldRule, newRule []string) error {
	return a.store.updateRow(ctx, ptype, oldRule, newRule)
}

func (a *Adapter) UpdatePoliciesCtx(ctx context.Context, sec string, ptype string, oldRules, newRules [][]string) error {
	return a.store.batchUpdate(ctx, ptype, oldRules, newRules)
}

func (a *Adapter) UpdateFilteredPoliciesCtx(ctx context.Context, sec string, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return nil, errors.New("not implemented")
}
