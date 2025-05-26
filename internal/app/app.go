package app

import (
	"context"
	"fmt"

	pgrepo "go-api-server/internal/adapter/repository/postgresql"
	appsvc "go-api-server/internal/app/service"
	"go-api-server/internal/config"
	domainrepo "go-api-server/internal/domain/repository"
	domainsvc "go-api-server/internal/domain/service"
	"go-api-server/internal/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CoreServices struct {
	PgPool      *pgxpool.Pool
	AppServices *Services
}

type Repositories struct {
	ConfigurationRepository domainrepo.ConfigurationRepository
}

type Services struct {
	ConfigurationService domainsvc.ConfigurationService
}

type CoreCleanupFunc func()

var log = logger.Instance()

func NewCoreServices(ctx context.Context, cfg *config.Config) (*CoreServices, CoreCleanupFunc, error) {
	log.Info("Initializing core services...")
	var cleanupFuncs []CoreCleanupFunc

	executeCleanups := func() {
		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}
	}

	pgPool, err := pgrepo.NewPostgreSQLPool(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("PostgreSQL connection failed: %w", err)
	}

	cleanupFuncs = append(cleanupFuncs, func() {
		pgPool.Close()
		log.Info("PostgreSQL connection pool closed")
	})

	core := &CoreServices{
		PgPool:      pgPool,
		AppServices: newAppServices(ctx, newRepositories(pgPool)),
	}

	finalCleanupFunc := func() {
		log.Info("Starting core services cleanup...")
		executeCleanups()
		log.Info("Core services cleanup finished")
	}

	log.Info("Core services initialized successfully")
	return core, finalCleanupFunc, nil
}

func newRepositories(pgPool *pgxpool.Pool) *Repositories {
	return &Repositories{
		ConfigurationRepository: pgrepo.NewConfigurationRepository(pgPool),
	}
}

func newAppServices(ctx context.Context, repositories *Repositories) *Services {
	return &Services{
		ConfigurationService: appsvc.NewConfigurationService(ctx, repositories.ConfigurationRepository),
	}
}
