package data

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/4ont/kratos-layout/internal/biz/sample"
	"github.com/4ont/kratos-layout/internal/data/model"
	"github.com/4ont/kratos-layout/internal/data/postgres"
	"github.com/4ont/kratos-layout/internal/data/redis"
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

func (a *AuthRepo) ValidAdminToken(ctx context.Context, token string) (*sample.Administrator, error) {
	tokenKey := fmt.Sprintf("%v%v", adminCacheKeyPrefix, token)
	strCmd := redis.GetRedis().Get(ctx, tokenKey)
	if strCmd.Err() != nil {
		return nil, errors.Wrap(strCmd.Err(), fmt.Sprintf("redis get %s failed", tokenKey))
	}
	var u sample.Administrator
	bys, err := strCmd.Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = json.Unmarshal(bys, &u); err != nil {
		return nil, errors.WithStack(err)
	}

	return &u, nil
}
