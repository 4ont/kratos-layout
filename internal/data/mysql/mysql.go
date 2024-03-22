package mysql

import (
	"math/rand"

	"github.com/Taskon-xyz/kit/go/kratostune"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	gormotel "gorm.io/plugin/opentelemetry/tracing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"

	"github.com/Taskon-xyz/kratos-layout/internal/conf"
)

var (
	gormConnections  map[string]*gorm.DB
	connections      map[string]*sqlx.DB
	connectionSlaves map[string][]*sqlx.DB
	slaveCounts      map[string]int
)

func InitDBConnections(dbCnf *conf.Data_Database) (func(), error) {
	gormConnections = make(map[string]*gorm.DB)
	connections = make(map[string]*sqlx.DB)
	connectionSlaves = make(map[string][]*sqlx.DB)
	// set logger
	gormCnf := &gorm.Config{}
	gormCnf.Logger = kratostune.NewGormLogger(log.With(kratostune.GetLogger(), "module", "gorm"))

	masterDB, err := gorm.Open(mysql.Open(dbCnf.GetDsn()), gormCnf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	slaves := make([]gorm.Dialector, 0, len(dbCnf.GetSlavesDsn()))
	for _, slave := range dbCnf.GetSlavesDsn() {
		slaves = append(slaves, mysql.Open(slave))
		dbx, err := sqlx.Open("mysql", slave)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		connectionSlaves[dbCnf.Name] = append(connectionSlaves[dbCnf.Name], dbx)
	}
	slaveCounts[dbCnf.Name] = len(dbCnf.GetSlavesDsn())

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

	gormConnections[dbCnf.Name] = masterDB
	connections[dbCnf.Name] = sqlx.NewDb(sqlDB, "mysql")

	cleanup := func() {
		for _, dbx := range connections {
			dbx.Close()
		}
		for _, slavers := range connectionSlaves {
			for _, db := range slavers {
				db.Close()
			}
		}
		gormConnections = nil
		connections = nil
		connectionSlaves = nil
	}

	return cleanup, nil
}

func GetDB(name ...string) *gorm.DB {
	if len(name) == 1 {
		return gormConnections[name[0]]
	}

	for _, db := range gormConnections {
		return db
	}

	return nil
}

func GetDBX(name ...string) *sqlx.DB {
	if len(name) == 1 {
		return connections[name[0]]
	}

	for _, db := range connections {
		return db
	}

	return nil
}

func GetDBXSlave(name ...string) *sqlx.DB {
	count := slaveCounts[name[0]]
	if len(name) == 1 {
		return connectionSlaves[name[0]][rand.Intn(count)]
	}

	for _, db := range connectionSlaves {
		return db[rand.Intn(count)]
	}

	return nil
}
