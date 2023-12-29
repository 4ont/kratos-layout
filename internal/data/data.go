package data

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/4ont/kratos-layout/internal/conf"
	"github.com/4ont/kratos-layout/internal/data/postgres"
	"github.com/4ont/kratos-layout/internal/data/redis"
)

// InitData .
func InitData(c *conf.Data) (func(), error) {
	cleanDB, err := postgres.InitDBConnections(c.GetDatabase())
	if err != nil {
		return nil, err
	}

	cleanRedis, err := redis.InitRedis(c.GetRedis())

	cleanup := func() {
		log.Info("closing the data resources")
		cleanDB()
		cleanRedis()
	}

	return cleanup, nil
}
