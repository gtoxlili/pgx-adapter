package pgxadapter

import (
	"context"
	"errors"
	"github.com/casbin/casbin/v2/model"
)

func (a *Adapter) LoadFilteredPolicyCtx(ctx context.Context, model model.Model, filter interface{}) error {
	return errors.New("not implemented")
}

func (a *Adapter) IsFilteredCtx(ctx context.Context) bool {
	return false
}
