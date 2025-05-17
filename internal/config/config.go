package config

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	RunMode string `mapstructure:"run_mode"`
}

var (
	v        *viper.Viper
	instance *Config
	once     sync.Once
	initErr  error
	mu       sync.RWMutex
)

func Init() error {
	once.Do(func() {
		v = viper.New()
		setViperDefaults(v)
		initializeReaderConfig(v)

		var cfg Config
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				slog.Info("Configuration file not found, using default values")
			} else {
				initErr = fmt.Errorf("failed to read configuration file: %w", err)
				return
			}
		}

		if err := v.Unmarshal(&cfg); err != nil {
			initErr = fmt.Errorf("failed to unmarshal initial configuration: %w", err)
			return
		}
		instance = &cfg

		configFileUsed := v.ConfigFileUsed()
		if configFileUsed != "" {
			slog.Info("Configuration successfully loaded and initialized from file", "file", configFileUsed)
		} else {
			slog.Info("Configuration successfully initialized using default values")
		}

		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			slog.Info("Configuration file change detected, reloading...", "file", e.Name)
			var newCfg Config
			if err := v.Unmarshal(&newCfg); err != nil {
				slog.Error("Failed to reload and unmarshal configuration file, continuing with previous settings", "file", e.Name, "error", err)
				return
			}

			mu.Lock()
			instance = &newCfg
			mu.Unlock()
			slog.Info("Configuration successfully reloaded", "file", e.Name)
		})
	})
	return initErr
}

func Instance() *Config {
	mu.RLock()
	if instance != nil {
		mu.RUnlock()
		return instance
	}
	mu.RUnlock()

	if err := Init(); err != nil {
		slog.Error("Failed to get configuration instance due to an error during initialization", "error", err)
		return nil
	}

	mu.RLock()
	defer mu.RUnlock()
	return instance
}

func setViperDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.run_mode", "debug")
}

func initializeReaderConfig(v *viper.Viper) {
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("internal/config")
	v.AddConfigPath(".")
}
