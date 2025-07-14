package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/teakingwang/ginmicro/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

type PGDBConfig struct {
	User     string
	Password string
	DBName   string
	DBSchema string
	Host     string
	Port     int
	Debug    bool
	LogLevel logger.LogLevel
}

func NewPostgres(c *config.DatabaseConfig) (*gorm.DB, error) {
	pDB, err := newPGDBWithLevel(
		c.User,
		c.Password,
		c.Database,
		c.Schema,
		c.Host,
		c.Level,
		c.Port,
	)
	if err != nil {
		logrus.Errorf("failed to connect database, %+v", err)
		return nil, err
	}

	return pDB, nil
}

func newPGDBWithLevel(user, password, db, schema, host, level string, port int) (*gorm.DB, error) {
	connURL := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s search_path=%s sslmode=disable",
		host, port, user, db, password, schema)

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

	pgdb, err := gorm.Open(postgres.Open(connURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	return pgdb, nil
}
