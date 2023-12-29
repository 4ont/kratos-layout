package server

import (
	"context"

	"github.com/4ont/kit/go/kratostune"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/4ont/kratos-layout/api/_pb/example"
	"github.com/4ont/kratos-layout/internal/biz/sample"
	"github.com/4ont/kratos-layout/internal/conf"
	"github.com/4ont/kratos-layout/internal/service"
)

// NewPortalHTTPServer new an HTTP server.
func NewPortalHTTPServer(c *conf.Server, admin *service.AdminService) *http.Server {
	var opts = []http.ServerOption{
		http.ResponseEncoder(kratostune.ResponseEncoder),
		http.ErrorEncoder(kratostune.ErrorEncoder(example.ErrorReason_value)),
	}
	if c.Portal.Network != "" {
		opts = append(opts, http.Network(c.Portal.Network))
	}
	if c.Portal.Addr != "" {
		opts = append(opts, http.Address(c.Portal.Addr))
	}
	if c.Portal.Timeout != nil {
		opts = append(opts, http.Timeout(c.Portal.Timeout.AsDuration()))
	}
	middlewares := kratostune.PrepareMiddleWare()
	middlewares = append(middlewares,
		NewAdminAuthMidware(),
	)
	opts = append(opts, http.Middleware(middlewares...))

	srv := http.NewServer(opts...)
	example.RegisterAdminHTTPServer(srv, admin)

	// 暴露监控采集端点
	srv.Handle("/metrics", promhttp.Handler())

	return srv
}

const (
	defaultAuthorizationHeader = "Authorization"
	adminAuthorizationMD       = "x-md-admin-token"
)

func NewAdminAuthMidware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var adminToken string
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				adminToken = md.Get(adminAuthorizationMD)[0]
			} else if tr, ok := transport.FromServerContext(ctx); ok {
				adminToken = tr.RequestHeader().Get(defaultAuthorizationHeader)
			} else {
				// 缺少可认证的token，返回错误
				return nil, errors.New("缺少可认证的token")
			}
			if adminToken == "" {
				return nil, errors.New("not login")
			}

			admin, err := sample.ValidAdminToken(ctx, adminToken)
			if err != nil {
				log.Error("ValidAdminToken failed", zap.Error(err))
				return nil, err
			}

			ctx = sample.ContextWithAdmin(ctx, admin)

			return handler(ctx, req)
		}
	}
}
