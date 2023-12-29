package sample

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware/selector"
)

// NewWhiteListMatcher jwt白名单，在名单中的路由不用校验jwt
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	// key的格式依据proto定义的: /package.service/rpcName
	whiteList["/probe.Probe/healthy"] = struct{}{}
	whiteList["/probe.Probe/ready"] = struct{}{}

	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
