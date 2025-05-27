package database

import (
	"context"
	"fmt"
	"time"

	"ditto/config"

	loggerCustom "ditto/pkg/logger"

	"github.com/spf13/cast"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Module = fx.Options(
	fx.Provide(NewDatabase),
)

type database struct {
	DB *gorm.DB
}

type Database interface {
	GetDB() *gorm.DB
	BeginTransaction() (*gorm.DB, error)
	ReleaseTransaction(tx *gorm.DB, err error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
}

func NewDatabase(lc fx.Lifecycle, config *config.Config, log loggerCustom.Logger) (Database, error) {
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				config.DB.Host, config.DB.User, config.DB.Password, config.DB.DBName, config.DB.Port, config.DB.SSLMode),
			PreferSimpleProtocol: true,
		},
	), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.WithError(err).Fatal("Failed to get DB from gorm")
		return nil, err
	}

	db.AutoMigrate()

	if config.DB.SetMaxIdleConns != "" {
		sqlDb.SetMaxIdleConns(cast.ToInt(config.DB.SetMaxIdleConns))
	}
	if config.DB.SetMaxOpenConns != "" {
		sqlDb.SetMaxOpenConns(cast.ToInt(config.DB.SetMaxOpenConns))
	}
	if config.DB.SetConnMaxLifetime != "" {
		sqlDb.SetConnMaxLifetime(cast.ToDuration(config.DB.SetConnMaxLifetime) * time.Hour)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sqlDb, err := db.DB()
			if err != nil {
				log.WithError(err).Fatal("Failed to get DB from gorm")
				return err
			}
			return sqlDb.Close()
		},
	})

	return &database{DB: db}, nil
}

// BeginTransaction implements Database.
func (d *database) BeginTransaction() (*gorm.DB, error) {
	return d.DB.Begin(), nil
}

// CommitTransaction implements Database.
func (d *database) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

// GetDB implements Database.
func (d *database) GetDB() *gorm.DB {
	return d.DB
}

// ReleaseTransaction implements Database.
func (d *database) ReleaseTransaction(tx *gorm.DB, err error) {
	if err != nil {
		d.RollbackTransaction(tx)
	}
	tx.Commit()
}

// RollbackTransaction implements Database.
func (d *database) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}
