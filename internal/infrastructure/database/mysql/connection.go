package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func LoadConfigFromEnv() (Config, error) {

	port, err := strconv.Atoi(getEnv("DB_PORT", "3306"))
	if err != nil {
		return Config{}, err
	}

	maxOpen, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return Config{}, err
	}

	maxIdle, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	if err != nil {
		return Config{}, err
	}

	lifetimeMinutes, err := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "5"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_NAME", "admin_videos"),

		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: time.Minute * time.Duration(lifetimeMinutes),
	}, nil
}

func NewConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
