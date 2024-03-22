package data

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/Taskon-xyz/kratos-layout/internal/conf"
	"github.com/Taskon-xyz/kratos-layout/internal/data/mysql"
	"github.com/Taskon-xyz/kratos-layout/internal/data/redis"
)

// InitData .
func InitData(c *conf.Data) (func(), error) {
	cleanDB, err := mysql.InitDBConnections(c.GetDatabase())
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
