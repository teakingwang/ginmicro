package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/teakingwang/ginmicro/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

type MySQLDBConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
	Debug    bool
	LogLevel logger.LogLevel
}

func NewMySQL(c *config.DatabaseConfig) (*gorm.DB, error) {
	mDB, err := newMySQLWithLevel(
		c.User,
		c.Password,
		c.Database,
		c.Host,
		c.Level,
		c.Port,
	)
	if err != nil {
		logrus.Errorf("failed to connect database, %+v", err)
		return nil, err
	}

	return mDB, nil
}

func newMySQLWithLevel(user, password, db, host, level string, port int) (*gorm.DB, error) {
	// MySQL DSN format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, db)

	var logLevel logger.LogLevel
	switch strings.ToLower(level) {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	default:
		logLevel = logger.Info
	}

	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
