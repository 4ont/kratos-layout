package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	_ "go.uber.org/automaxprocs"

	"github.com/4ont/kit/go/kratostune"

	"github.com/4ont/kratos-layout/internal/biz/sample"
	"github.com/4ont/kratos-layout/internal/conf"
	"github.com/4ont/kratos-layout/internal/data"
	"github.com/4ont/kratos-layout/internal/global"
	"github.com/4ont/kratos-layout/internal/server"
	"github.com/4ont/kratos-layout/internal/service"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
	// flagconfsrc is the config source flag.
	flagconfsrc string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&flagconfsrc, "conf-src", "file", "config source, eg: -conf-src aws-appconfig")
	flag.Usage = usage
}

func newApp(confServer *conf.Server) *kratos.App {
	sample.RegisterRepo(data.NewAuthRepo())
	probeService := service.NewProbeService()
	authService := service.NewSampleService()
	adminService := service.NewAdminService()

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			server.NewGRPCServer(confServer, probeService, authService),
			server.NewHTTPServer(confServer, probeService, authService),
			server.NewPortalHTTPServer(confServer, adminService),
		),
	)
}

func main() {
	flag.Parse()

	// 加载配置
	cleanconf := global.InitConfig(flagconfsrc, flagconf)
	defer cleanconf()
	bc := global.GetConfig()
	os.Setenv("ENV", bc.Env)

	// 初始化日志
	kratostune.InitLogger(Name, Version, bc.LogLevel)
	// 初始化分布式追踪
	err := kratostune.SetTracerProvider(bc.Tracing.GetType(), bc.Tracing.GetHost(), int(bc.Tracing.GetPort()),
		Name, Version, bc.Env)
	if err != nil {
		log.Errorf("SetTracerProvider failed: %+v", err)
		panic(err)
	}
	log.Infof("SetTracerProvider, type: %s, host: %s, port: %d", bc.Tracing.GetType(), bc.Tracing.GetHost(), int(bc.Tracing.GetPort()))

	//  初始化 sentry
	err = kratostune.InitSentry(Name, Version, bc.Env, bc.Sentry.GetDsn(), bc.Sentry.GetAttachStackTrace())
	if err != nil {
		log.Errorf("InitSentry: %+v", err)
		panic(err)
	}

	//  初始化服务
	cleanup, err := data.InitData(bc.Data)
	if err != nil {
		log.Errorf("wire up failed: %+v", err)
		panic(err)
	}
	defer cleanup()

	app := newApp(bc.Server)
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s Version: %s\n", Name, Version)
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
