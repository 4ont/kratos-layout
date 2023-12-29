package server

import (
	"github.com/4ont/kit/go/kratostune"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/4ont/kratos-layout/api/_pb/example"
	probeapi "github.com/4ont/kratos-layout/api/_pb/probe"
	bizauth "github.com/4ont/kratos-layout/internal/biz/sample"
	"github.com/4ont/kratos-layout/internal/conf"
	"github.com/4ont/kratos-layout/internal/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	probe *service.ProbeService,
	auth *service.SampleService) *http.Server {
	var opts = []http.ServerOption{
		http.ResponseEncoder(kratostune.ResponseEncoder),
		http.ErrorEncoder(kratostune.ErrorEncoder(example.ErrorReason_value)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	middlewares := kratostune.PrepareMiddleWare()
	middlewares = append(middlewares,
		selector.Server( //jwt中间件
			kratostune.NewJWTServerMidware(kratostune.NewJWT(c.GetAuth().GetJwtKey(), 0)),
		).Match(bizauth.NewWhiteListMatcher()).Build(),
	)
	opts = append(opts, http.Middleware(middlewares...))

	srv := http.NewServer(opts...)
	probeapi.RegisterProbeHTTPServer(srv, probe)
	example.RegisterSampleHTTPServer(srv, auth)

	// 暴露监控采集端点
	srv.Handle("/metrics", promhttp.Handler())

	return srv
}
