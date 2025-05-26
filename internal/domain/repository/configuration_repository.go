package repository

import (
	"context"

	"go-api-server/internal/adapter/repository/postgresql/schema"
)

type ConfigurationRepository interface {
	FindByKey(ctx context.Context, key string) (*schema.FindByKeySchema, error)
}
