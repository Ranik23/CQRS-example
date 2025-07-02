package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type StorageConfig struct {
	Main PostgresConfig `mapstructure:"main"`
	Side PostgresConfig `mapstructure:"side"`
}

type PostgresConfig struct {
	Host               string     `mapstructure:"host"`
	Port               string     `mapstructure:"port"`
	Username           string     `mapstructure:"username"`
	Password           string     `mapstructure:"password"`
	Database           string     `mapstructure:"dbname"`
	ConnectionAttempts int        `mapstructure:"connection_attempts"`
	PoolConfig         PoolConfig `mapstructure:"pool"`
	OutboxTable        OutboxConfig `mapstructure:"outbox_table"`
}

type PoolConfig struct {
	MaxConnections    int `mapstructure:"max_connections"`
	MinConnections    int `mapstructure:"min_connections"`
	MaxLifeTime       int `mapstructure:"max_lifetime"`
	MaxIdleTime       int `mapstructure:"max_idle_time"`
	HealthCheckPeriod int `mapstructure:"health_check_period"`
}

type OutboxConfig struct {
	BatchSize   int `mapstructure:"batch_size"`
	NumWorkers  int `mapstructure:"num_workers"`
}

func (s *PostgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", s.Host, s.Port, s.Username, s.Password, s.Database)
}

func (s *PostgresConfig) ApplyMigrations(ctx context.Context, migrationsPath string) error {
	dsn := s.GetDSN()

	fmt.Println("Applying migrations to ", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (s *PostgresConfig) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(s.GetDSN())
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(s.PoolConfig.MaxConnections)
	poolConfig.MinConns = int32(s.PoolConfig.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(s.PoolConfig.MaxLifeTime) * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(s.PoolConfig.MaxIdleTime) * time.Second
	poolConfig.HealthCheckPeriod = time.Duration(s.PoolConfig.HealthCheckPeriod) * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}


	fmt.Println("Successfully created connection pool")

	return pool, nil
}
