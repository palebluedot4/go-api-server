package service

import (
	"context"
	"errors"
	"fmt"

	"go-api-server/internal/domain/model"
	"go-api-server/internal/domain/repository"
	"go-api-server/internal/domain/service"

	"github.com/jackc/pgx/v5"
)

type configurationServiceImpl struct {
	configurationRepository repository.ConfigurationRepository
}

func NewConfigurationService(ctx context.Context, configurationRepository repository.ConfigurationRepository) service.ConfigurationService {
	return &configurationServiceImpl{
		configurationRepository: configurationRepository,
	}
}

func (s *configurationServiceImpl) GetConfigurationByKey(ctx context.Context, key string) (*model.GetConfigurationByKey, error) {
	configuration, err := s.configurationRepository.FindByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("configuration with key %s not found", key)
		}
		return nil, fmt.Errorf("failed to get configuration by key %s: %w", key, err)
	}

	if configuration == nil {
		return nil, fmt.Errorf("configuration with key %s not found", key)
	}

	return &model.GetConfigurationByKey{
		ConfigValue: configuration.ConfigValue,
	}, nil
}
