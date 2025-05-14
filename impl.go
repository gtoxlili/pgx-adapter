package pgxadapter

import (
	"context"
	"github.com/casbin/casbin/v2/model"
)

func (a *Adapter) LoadPolicyCtx(ctx context.Context, model model.Model) error {
	//TODO implement me
	panic("implement me")
}

func (a *Adapter) SavePolicyCtx(ctx context.Context, model model.Model) error {
	//TODO implement me
	panic("implement me")
}

func (a *Adapter) AddPolicyCtx(ctx context.Context, sec string, ptype string, rule []string) error {
	//TODO implement me
	panic("implement me")
}

func (a *Adapter) RemovePolicyCtx(ctx context.Context, sec string, ptype string, rule []string) error {
	//TODO implement me
	panic("implement me")
}

func (a *Adapter) RemoveFilteredPolicyCtx(ctx context.Context, sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	//TODO implement me
	panic("implement me")
}
