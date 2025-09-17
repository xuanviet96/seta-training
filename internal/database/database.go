package database

import (
	"database/sql"
	"errors"

	"github.com/xuanviet96/seta-training/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string, log *zap.Logger) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is empty")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Initialize opens a DB connection using the provided config.
func Initialize(cfg config.Config) (*gorm.DB, error) {
	return Connect(cfg.DatabaseURL, nil)
}

func Ping(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func RawPing(sqlDB *sql.DB) error {
	return sqlDB.Ping()
}
