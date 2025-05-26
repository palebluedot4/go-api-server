package service_test

import (
	"context"
	"errors"
	"testing"

	"go-api-server/internal/adapter/repository/postgresql/schema"
	"go-api-server/internal/app/service"
	"go-api-server/internal/domain/model"
	"go-api-server/internal/domain/repository/mocks"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationService_GetConfigurationByKey(t *testing.T) {
	ctx := context.Background()

	mockRepository := mocks.NewMockConfigurationRepository(t)
	configService := service.NewConfigurationService(ctx, mockRepository)

	testKey := "app_version"
	testValue := "1.0.0"

	t.Run("success", func(t *testing.T) {
		expectedSchema := &schema.FindByKeySchema{
			ConfigValue: testValue,
		}
		expectedModel := &model.GetConfigurationByKey{
			ConfigValue: expectedSchema.ConfigValue,
		}

		mockRepository.On("FindByKey", ctx, testKey).Return(expectedSchema, nil).Once()

		model, err := configService.GetConfigurationByKey(ctx, testKey)

		assert.NoError(t, err)
		assert.NotNil(t, model)
		assert.Equal(t, expectedModel.ConfigValue, model.ConfigValue)
	})

	t.Run("repository error other than not found", func(t *testing.T) {
		expectedError := errors.New("database error")
		mockRepository.On("FindByKey", ctx, testKey).Return(nil, expectedError).Once()

		model, err := configService.GetConfigurationByKey(ctx, testKey)

		assert.Error(t, err)
		assert.Nil(t, model)
		assert.EqualError(t, err, "failed to get configuration by key "+testKey+": "+expectedError.Error())
	})

	t.Run("not found from repository (pgx.ErrNoRows)", func(t *testing.T) {
		mockRepository.On("FindByKey", ctx, testKey).Return(nil, pgx.ErrNoRows).Once()

		model, err := configService.GetConfigurationByKey(ctx, testKey)

		assert.Error(t, err)
		assert.Nil(t, model)
		assert.EqualError(t, err, "configuration with key "+testKey+" not found")
	})

	t.Run("not found when repository returns nil, nil (should be handled by service)", func(t *testing.T) {
		mockRepository.On("FindByKey", ctx, testKey).Return(nil, nil).Once()

		model, err := configService.GetConfigurationByKey(ctx, testKey)

		assert.Error(t, err)
		assert.Nil(t, model)
		assert.EqualError(t, err, "configuration with key "+testKey+" not found")
	})
}
