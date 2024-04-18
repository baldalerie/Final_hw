package database

import (
	"FinalTaskAppGoBasic/internal/configs"
	"FinalTaskAppGoBasic/internal/logs"
	"FinalTaskAppGoBasic/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	db *gorm.DB
}

func (c *Connection) Gorm() *gorm.DB {
	return c.db
}

func (c *Connection) Connect(cfg *configs.Database) error {
	logEntry := logs.Log.
		WithField("host", cfg.Host).
		WithField("port", cfg.Port).
		WithField("user", cfg.User).
		WithField("db_name", cfg.DBName).
		WithField("ssl_mode", cfg.SSLMode).
		WithField("time_zone", cfg.TimeZone)

	logEntry.Info("connecting to database")

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	var err error
	c.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logEntry.WithError(err).Error("connect to database failed")

		return err
	}

	logEntry.Info("connect to database succeed")

	return nil
}

func (c *Connection) Migrate() error {
	err := c.db.AutoMigrate(&models.Users{})
	if err != nil {
		logs.Log.WithError(err).Error("migrate users failed")

		return err
	}

	logs.Log.Info("migrate users succeed")

	err = c.db.AutoMigrate(&models.Transactions{})
	if err != nil {
		logs.Log.WithError(err).Error("migrate transactions failed")

		return err
	}

	logs.Log.Info("migrate transactions succeed")

	return nil
}

func (c *Connection) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		logs.Log.WithError(err).Error("get database connection failed")

		return err
	}

	err = sqlDB.Close()
	if err != nil {
		logs.Log.WithError(err).Error("close database connection failed")

		return err
	}

	logs.Log.Info("close database connection succeed")

	return nil
}

func New() *Connection {
	return &Connection{}
}
