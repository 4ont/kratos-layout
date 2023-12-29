package sample

import (
	"context"
)

type Administrator struct {
	UserID         string `json:"user_id"`
	WalletAddr     string `json:"wallet_addr"`
	UniqueNickname string `json:"unique_nickname"`
	Avatar         string `json:"avatar"`
	Enabled        bool   `json:"enabled"`
	IsAdmin        bool   `json:"is_admin"`
	PartnerID      string `json:"partner_id"`
	AppID          string `json:"app_id"`
	IsSuperAdmin   bool   `json:"is_super_admin"`
}

type adminContextKey struct{}

var (
	defaultAdminContextKey = adminContextKey{}
)

func ValidAdminToken(ctx context.Context, token string) (*Administrator, error) {
	return authReop.ValidAdminToken(ctx, token)
}

func ContextWithAdmin(ctx context.Context, user *Administrator) context.Context {
	newCtx := context.WithValue(ctx, defaultAdminContextKey, user)
	return newCtx
}

func AdminFromContext(ctx context.Context) (*Administrator, bool) {
	u, ok := ctx.Value(defaultAdminContextKey).(*Administrator)
	return u, ok
}
