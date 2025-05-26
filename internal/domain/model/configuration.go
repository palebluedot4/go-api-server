package model

import (
	"time"
)

type ApplicationConfiguration struct {
	ConfigKey   string    `json:"config_key" `
	ConfigValue string    `json:"config_value" `
	ValueType   string    `json:"value_type" `
	Description *string   `json:"description,omitempty" `
	IsActive    bool      `json:"is_active" `
	CreatedAt   time.Time `json:"created_at" `
	UpdatedAt   time.Time `json:"updated_at" `
}

type GetConfigurationByKey struct {
	ConfigValue string `json:"config_value"`
}
