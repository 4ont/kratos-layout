package data

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Taskon-xyz/kratos-layout/internal/biz/sample"
	"github.com/Taskon-xyz/kratos-layout/internal/data/model"
	"github.com/Taskon-xyz/kratos-layout/internal/data/mysql"
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
	pgdb := mysql.GetDB().WithContext(ctx)
	user, err := model.User{}.SelectOneByEmail(pgdb, email)
	if err != nil {
		errors.Wrap(err, "query user by email")
		return
	}
	return user != nil, nil
}
