package database

import (
	"context"

	"gorm.io/gorm"
)

// ConnectFromEnv opens a Postgres connection using environment variables.
func ConnectFromEnv(ctx context.Context) (*gorm.DB, error) {
	cfg := ConfigFromEnv()
	return Connect(ctx, cfg)
}

// Connect opens a Postgres connection and validates it.
func Connect(ctx context.Context, cfg PostgresConfig) (*gorm.DB, error) {
	db, err := OpenPostgres(cfg)
	if err != nil {
		return nil, err
	}

	if err := AutoMigrate(db); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
