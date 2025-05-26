package postgresql

import (
	"context"

	"go-api-server/internal/adapter/repository/postgresql/schema"
	"go-api-server/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type configurationRepository struct {
	db *pgxpool.Pool
}

func NewConfigurationRepository(db *pgxpool.Pool) repository.ConfigurationRepository {
	return &configurationRepository{
		db: db,
	}
}

func (r *configurationRepository) FindByKey(ctx context.Context, key string) (*schema.FindByKeySchema, error) {
	query := `
		SELECT
			config_value
		FROM
			application_configurations
		WHERE
			config_key = $1 AND
			is_active
	`

	var config schema.FindByKeySchema
	err := r.db.QueryRow(ctx, query, key).Scan(&config.ConfigValue)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
