package postgresql

import (
	"context"
	"fmt"

	"go-api-server/internal/config"
	"go-api-server/internal/pkg/logger"
	"go-api-server/internal/pkg/timeouts"

	"github.com/jackc/pgx/v5/pgxpool"
)

var log = logger.Instance()

func NewPostgreSQLPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	if !cfg.Storage.PostgreSQL.Enabled {
		log.Info("PostgreSQL is disabled in configuration")
		return nil, nil
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Storage.PostgreSQL.Host,
		cfg.Storage.PostgreSQL.Port,
		cfg.Storage.PostgreSQL.User,
		cfg.Storage.PostgreSQL.Password,
		cfg.Storage.PostgreSQL.DBName,
		cfg.Storage.PostgreSQL.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL DSN: %w", err)
	}

	if cfg.Storage.PostgreSQL.MaxOpenConns > 0 {
		poolConfig.MaxConns = int32(cfg.Storage.PostgreSQL.MaxOpenConns)
	}
	if cfg.Storage.PostgreSQL.MaxIdleConns > 0 {
		poolConfig.MinConns = int32(cfg.Storage.PostgreSQL.MaxIdleConns)
	}
	if cfg.Storage.PostgreSQL.ConnMaxLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.Storage.PostgreSQL.ConnMaxLifetime
	}
	if cfg.Storage.PostgreSQL.ConnMaxIdleTime > 0 {
		poolConfig.MaxConnIdleTime = cfg.Storage.PostgreSQL.ConnMaxIdleTime
	}

	connectTimeout := timeouts.StorageConnect(cfg)
	poolCtx, poolCancel := context.WithTimeout(ctx, connectTimeout)
	defer poolCancel()

	pool, err := pgxpool.NewWithConfig(poolCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL connection pool: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, connectTimeout)
	defer pingCancel()
	if err = pool.Ping(pingCtx); err != nil {
		log.WithField("error", err).WithField("host", cfg.Storage.PostgreSQL.Host).WithField("database", cfg.Storage.PostgreSQL.DBName).Error("Failed to ping PostgreSQL database")
		pool.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL database within %v: %w", connectTimeout, err)
	}

	log.Infof("Successfully connected to PostgreSQL database: %s on host: %s", cfg.Storage.PostgreSQL.DBName, cfg.Storage.PostgreSQL.Host)
	return pool, nil
}
