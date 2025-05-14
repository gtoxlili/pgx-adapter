package pgxadapter

import (
	"context"
	"github.com/casbin/casbin/v2/model"
)

func (a *Adapter) LoadPolicy(model model.Model) error {
	return a.LoadPolicyCtx(context.Background(), model)
}

func (a *Adapter) SavePolicy(model model.Model) error {
	return a.SavePolicyCtx(context.Background(), model)
}

func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return a.AddPolicyCtx(context.Background(), sec, ptype, rule)
}

func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return a.RemovePolicyCtx(context.Background(), sec, ptype, rule)
}

func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return a.RemoveFilteredPolicyCtx(context.Background(), sec, ptype, fieldIndex, fieldValues...)
}

func (a *Adapter) UpdatePolicy(sec string, ptype string, oldRule, newRule []string) error {
	return a.UpdatePolicyCtx(context.Background(), sec, ptype, oldRule, newRule)
}

func (a *Adapter) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error {
	return a.UpdatePoliciesCtx(context.Background(), sec, ptype, oldRules, newRules)
}

func (a *Adapter) UpdateFilteredPolicies(sec string, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return a.UpdateFilteredPoliciesCtx(context.Background(), sec, ptype, newRules, fieldIndex, fieldValues...)
}

func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	return a.AddPoliciesCtx(context.Background(), sec, ptype, rules)
}

func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return a.RemovePoliciesCtx(context.Background(), sec, ptype, rules)
}

func (a *Adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	return a.LoadFilteredPolicyCtx(context.Background(), model, filter)
}

func (a *Adapter) IsFiltered() bool {
	return a.IsFilteredCtx(context.Background())
}
