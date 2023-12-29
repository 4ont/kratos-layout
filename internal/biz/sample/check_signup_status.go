package sample

import (
	"context"

	"github.com/pkg/errors"
)

const (
	checkSignupStatusTypeEmail = checkSignupStatusType("email")
)

var (
	authReop AuthRepo
)

type AuthRepo interface {
	CheckEmailSignupStatus(ctx context.Context, email string) (registered bool, err error)
	ValidAdminToken(ctx context.Context, token string) (*Administrator, error)
}

type checkSignupStatusType string

func RegisterRepo(repo AuthRepo) {
	authReop = repo
}

func CheckSignupStatus(ctx context.Context, signupType string, value string) (registered bool, err error) {
	if value == "" {
		return false, errors.New("invalid parameters")
	}
	switch checkSignupStatusType(signupType) {
	case checkSignupStatusTypeEmail:
		return authReop.CheckEmailSignupStatus(ctx, value)
	default:
		return false, errors.New("invalid parameters")
	}
}
