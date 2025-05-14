package pgxadapter

import (
	"context"
	"github.com/casbin/casbin/v2/model"
)

func (a *Adapter) LoadFilteredPolicyCtx(ctx context.Context, model model.Model, filter interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (a *Adapter) IsFilteredCtx(ctx context.Context) bool {
	//TODO implement me
	panic("implement me")
}
