package server

import (
	"github.com/Taskon-xyz/kit/go/kratostune"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/Taskon-xyz/kratos-layout/api/_pb/example"
	probeapi "github.com/Taskon-xyz/kratos-layout/api/_pb/probe"
	"github.com/Taskon-xyz/kratos-layout/internal/conf"
	"github.com/Taskon-xyz/kratos-layout/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server,
	probe *service.ProbeService,
	auth *service.SampleService) *grpc.Server {
	var opts []grpc.ServerOption
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	middlewares := kratostune.PrepareMiddleWare()
	opts = append(opts, grpc.Middleware(middlewares...))

	srv := grpc.NewServer(opts...)
	probeapi.RegisterProbeServer(srv, probe)
	example.RegisterSampleServer(srv, auth)
	return srv
}
