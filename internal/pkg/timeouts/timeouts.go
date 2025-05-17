package timeouts

import (
	"time"

	"go-api-server/internal/config"
)

const (
	defaultServerShutdownTimeout  = 30 * time.Second
	defaultServerReadTimeout      = 15 * time.Second
	defaultServerWriteTimeout     = 15 * time.Second
	defaultServerIdleTimeout      = 60 * time.Second
	defaultStorageConnectTimeout  = 5 * time.Second
	defaultStorageShutdownTimeout = 15 * time.Second
)

func ServerShutdown(cfg *config.Config) time.Duration {
	if cfg.Server.ShutdownTimeout == 0 {
		return defaultServerShutdownTimeout
	}
	return cfg.Server.ShutdownTimeout
}

func ServerRead(cfg *config.Config) time.Duration {
	if cfg.Server.ReadTimeout == 0 {
		return defaultServerReadTimeout
	}
	return cfg.Server.ReadTimeout
}

func ServerWrite(cfg *config.Config) time.Duration {
	if cfg.Server.WriteTimeout == 0 {
		return defaultServerWriteTimeout
	}
	return cfg.Server.WriteTimeout
}

func ServerIdle(cfg *config.Config) time.Duration {
	if cfg.Server.IdleTimeout == 0 {
		return defaultServerIdleTimeout
	}
	return cfg.Server.IdleTimeout
}

func StorageConnect(cfg *config.Config) time.Duration {
	if cfg.Storage.ConnectTimeout == 0 {
		return defaultStorageConnectTimeout
	}
	return cfg.Storage.ConnectTimeout
}

func StorageShutdown(cfg *config.Config) time.Duration {
	if cfg.Storage.ShutdownTimeout == 0 {
		return defaultStorageShutdownTimeout
	}
	return cfg.Storage.ShutdownTimeout
}
