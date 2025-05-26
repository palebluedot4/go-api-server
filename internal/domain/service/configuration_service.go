package service

import (
	"context"

	"go-api-server/internal/domain/model"
)

type ConfigurationService interface {
	GetConfigurationByKey(ctx context.Context, key string) (*model.GetConfigurationByKey, error)
}
