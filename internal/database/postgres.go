package database

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"main/internal/repository/model"
)

// PostgresConfig defines connection settings for Postgres.
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

// ConfigFromEnv builds a PostgresConfig using environment variables.
func ConfigFromEnv() PostgresConfig {
	return PostgresConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5452),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "9843"),
		Name:     getEnv("DB_NAME", "sandbox"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		TimeZone: getEnv("DB_TIMEZONE", "UTC"),
	}
}

// DSN returns the GORM Postgres DSN string.
func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
		c.TimeZone,
	)
}

// OpenPostgres opens a GORM connection to Postgres.
func OpenPostgres(cfg PostgresConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
}

// AutoMigrate runs the schema migrations for persistence models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Sandbox{}, &model.User{}, &model.RefreshToken{})
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return parsed
}
