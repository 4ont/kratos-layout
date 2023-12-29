package server

import (
	"github.com/4ont/kit/go/kratostune"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/4ont/kratos-layout/api/_pb/example"
	"github.com/4ont/kratos-layout/internal/conf"
	"github.com/4ont/kratos-layout/internal/service"
)

// NewPortalHTTPServer new an HTTP server.
func NewPortalHTTPServer(c *conf.Server, portalSvc *service.PortalService) *http.Server {
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
	opts = append(opts, http.Middleware(middlewares...))

	srv := http.NewServer(opts...)
	example.RegisterPortalHTTPServer(srv, portalSvc)

	// 暴露监控采集端点
	srv.Handle("/metrics", promhttp.Handler())

	return srv
}
