package global

import (
	"sync"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"

	"github.com/Taskon-xyz/kratos-layout/internal/conf"
)

var (
	cnf  *conf.Bootstrap
	lock sync.Mutex
)

func SetConfig(c *conf.Bootstrap) {
	lock.Lock()
	defer lock.Unlock()
	cnf = c
}

func GetConfig() *conf.Bootstrap {
	lock.Lock()
	defer lock.Unlock()

	return cnf
}

func InitConfig(confPath string) (clean func()) {
	src := file.NewSource(confPath)
	// 加载配置
	c := config.New(
		config.WithSource(src),
	)
	clean = func() { _ = c.Close() }

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	SetConfig(&bc)

	return
}

//---
