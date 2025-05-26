package schema

import (
	"database/sql"
	"time"
)

type ApplicationConfiguration struct {
	ConfigKey   string         `db:"config_key" `
	ConfigValue string         `db:"config_value" `
	ValueType   string         `db:"value_type" `
	Description sql.NullString `db:"description,omitempty" `
	IsActive    bool           `db:"is_active" `
	CreatedAt   time.Time      `db:"created_at" `
	UpdatedAt   time.Time      `db:"updated_at" `
}

type FindByKeySchema struct {
	ConfigValue string `db:"config_value"`
}
