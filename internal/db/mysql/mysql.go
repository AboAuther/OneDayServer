package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"one-day-server/configs"
)

var (
	db           *gorm.DB
	CustomizedDB *gorm.DB
)

const (
	dbConfig = "?parseTime=true&interpolateParams=true"
)

func Init() {
	var err error
	newLogger := newLogger()

	dsn := configs.GetEnvDefault("MYSQL_DSN", "")
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.MustGetEnv("MYSQL_USER"), configs.MustGetEnv("MYSQL_PASSWORD"), configs.MustGetEnv("MYSQL_HOST"), configs.MustGetEnv("MYSQL_PORT"), configs.MustGetEnv("MYSQL_DBNAME"))
	}
	dsn = dsn + dbConfig
	logrus.Infof("mysql: %s:password@tcp(%s:%s)/%s", configs.MustGetEnv("MYSQL_USER"), configs.MustGetEnv("MYSQL_HOST"), configs.MustGetEnv("MYSQL_PORT"), configs.MustGetEnv("MYSQL_DBNAME"))

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	logrus.Infof("init mysql db connection.")
}

func DB() *gorm.DB {
	if db == nil {
		Init()
	}
	return db
}

func Guard(tx *gorm.DB) {
	if tx.Error != nil {
		panic(tx.Error)
	}
}

func IsMissing(tx *gorm.DB) bool {
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return true
	} else if tx.Error == nil {
		return false
	} else {
		panic(tx.Error)
	}
}

func newLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
		},
	)

}
func InitCustomizedDB(dbName string) {
	var err error
	newLogger := newLogger()

	dsn := fmt.Sprintf("%s/%s%s", os.Getenv("MYSQL_DSN_BASE"), dbName, dbConfig)
	logrus.Infof("mysql: %s", dsn)

	CustomizedDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := CustomizedDB.DB()
	if err != nil {
		logrus.Panicf("get snapshot db failed, err: %s", err)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("init mysql db: %s connection.", dbName)
}
