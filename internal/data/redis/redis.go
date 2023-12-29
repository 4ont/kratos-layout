package redis

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"github.com/4ont/kratos-layout/internal/conf"
)

// common variables
const (
	defaultPoolSize    = 60
	defaultMinIdleConn = 20
	defaultMaxRetries  = 3
	defaultIdleTimeout = 5 * time.Minute
)

var (
	clusterClient *redis.ClusterClient
	redisClient   *redis.Client
)

func InitRedis(c *conf.Data_Redis) (func(), error) {
	if c.IsCluster {
		return initRedisCluster(c)
	} else {
		return initRedis(c)
	}
}

func initRedisCluster(c *conf.Data_Redis) (func(), error) {
	poolSize := c.GetPoolSize()
	if poolSize == 0 {
		poolSize = defaultPoolSize
	}
	minIdleConn := c.MinIdleConn
	if minIdleConn <= 0 {
		minIdleConn = defaultMinIdleConn
	}

	clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           c.GetAddrs(),
		Username:        c.GetUsername(),
		Password:        c.GetPassword(),
		MaxRetries:      defaultMaxRetries,
		ReadTimeout:     c.GetReadTimeout().AsDuration(),
		WriteTimeout:    c.GetWriteTimeout().AsDuration(),
		PoolSize:        int(poolSize),
		MinIdleConns:    int(minIdleConn),
		ConnMaxIdleTime: defaultIdleTimeout,
	})
	status := clusterClient.Ping(context.Background())
	if status.Err() != nil {
		return nil, errors.WithStack(status.Err())
	}

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(clusterClient); err != nil {
		return nil, err
	}
	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(clusterClient); err != nil {
		return nil, err
	}

	cleanup := func() {
		clusterClient.Close()
		log.Info("redis closed")
	}

	return cleanup, nil
}

func initRedis(c *conf.Data_Redis) (func(), error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     c.Addrs[0],
		Username: c.Username,
		Password: c.Password, // no password set
		DB:       int(c.Db),  // use default DB
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return nil, err
	}
	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(redisClient); err != nil {
		return nil, err
	}

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		redisClient.Close()
		log.Info("redis closed")
	}

	return cleanup, nil
}

func GetRedis() redis.UniversalClient {
	if clusterClient != nil {
		return clusterClient
	} else {
		return redisClient
	}
}
