package data

import (
	"context"

	"github.com/pkg/errors"

	"github.com/4ont/kratos-layout/internal/biz/sample"
	"github.com/4ont/kratos-layout/internal/data/model"
	"github.com/4ont/kratos-layout/internal/data/postgres"
)

const (
	adminCacheKeyPrefix = "admin_token:"
)

type AuthRepo struct {
}

func NewAuthRepo() sample.AuthRepo {
	return &AuthRepo{}
}

func (a *AuthRepo) CheckEmailSignupStatus(ctx context.Context, email string) (registered bool, err error) {
	pgdb := postgres.GetDB().WithContext(ctx)
	user, err := model.User{}.SelectOneByEmail(pgdb, email)
	if err != nil {
		errors.Wrap(err, "query user by email")
		return
	}
	return user != nil, nil
}
