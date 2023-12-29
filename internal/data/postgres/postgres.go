package postgres

import (
	"github.com/4ont/kit/go/kratostune"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	gormotel "gorm.io/plugin/opentelemetry/tracing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"

	"github.com/4ont/kratos-layout/internal/conf"
)

var (
	connections map[string]*gorm.DB
)

func InitDBConnections(dbCnf *conf.Data_Database) (func(), error) {
	connections = make(map[string]*gorm.DB)
	// set logger
	gormCnf := &gorm.Config{}
	gormCnf.Logger = kratostune.NewGormLogger(log.With(kratostune.GetLogger(), "module", "gorm"))

	masterDB, err := gorm.Open(postgres.Open(dbCnf.GetDsn()), gormCnf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	slaves := make([]gorm.Dialector, 0, len(dbCnf.GetSlavesDsn()))
	for _, slave := range dbCnf.GetSlavesDsn() {
		slaves = append(slaves, postgres.Open(slave))
	}

	// 设置主从
	if len(slaves) > 0 {
		err = masterDB.Use(dbresolver.Register(dbresolver.Config{
			Replicas: slaves,
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// 设置分布式追踪和监控采集
	err = masterDB.Use(gormotel.NewPlugin())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sqlDB, err := masterDB.DB()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// configure pool parameters
	sqlDB.SetMaxOpenConns(int(dbCnf.MaxOpenConnections))
	sqlDB.SetMaxIdleConns(int(dbCnf.MaxIdleConnections))
	sqlDB.SetConnMaxLifetime(dbCnf.ConnMaxLifetime.AsDuration())
	sqlDB.SetConnMaxIdleTime(dbCnf.ConnMaxIdleTime.AsDuration())

	connections[dbCnf.Name] = masterDB

	cleanup := func() {
		for k, _ := range connections {
			delete(connections, k)
		}
	}

	return cleanup, nil
}

func GetDB() *gorm.DB {
	if len(connections) == 1 {
		for _, db := range connections {
			return db
		}
	}

	return connections["default"]
}
