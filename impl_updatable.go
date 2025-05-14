package pgxadapter

import "context"

func (a *Adapter) UpdatePolicyCtx(ctx context.Context, sec string, ptype string, oldRule, newRule []string) error {
	//TODO implement me
	panic("implement me")
}

func (a *Adapter) UpdatePoliciesCtx(ctx context.Context, sec string, ptype string, oldRules, newRules [][]string) error {
	//TODO implement me
	panic("implement me")
}

func (a *Adapter) UpdateFilteredPoliciesCtx(ctx context.Context, sec string, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	//TODO implement me
	panic("implement me")
}
