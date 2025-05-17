package config

import (
	"fmt"
	"sync"
	"time"

	"go-api-server/internal/pkg/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Storage StorageConfig `mapstructure:"storage"`
}

type ServerConfig struct {
	Port            int           `mapstructure:"port"`
	RunMode         string        `mapstructure:"run_mode"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	Logger          logger.Config `mapstructure:"logger"`
}

type StorageConfig struct {
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

var (
	v        *viper.Viper
	instance *Config
	once     sync.Once
	initErr  error
	mu       sync.RWMutex
)

var log = logger.Instance()

func Init() error {
	once.Do(func() {
		v = viper.New()
		setViperDefaults(v)
		initializeReaderConfig(v)

		var cfg Config
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Info("Configuration file not found, using default values")
			} else {
				initErr = fmt.Errorf("failed to read configuration file: %w", err)
				log.WithField("error", err).Error("Failed to read configuration file")
				return
			}
		}

		if err := v.Unmarshal(&cfg); err != nil {
			initErr = fmt.Errorf("failed to unmarshal initial configuration: %w", err)
			log.WithField("error", err).Error("Failed to unmarshal initial configuration")
			return
		}
		instance = &cfg

		configFileUsed := v.ConfigFileUsed()
		if configFileUsed != "" {
			log.WithField("file", configFileUsed).Info("Configuration successfully loaded and initialized from file")
		} else {
			log.Info("Configuration successfully initialized using default values")
		}

		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.WithField("file", e.Name).Info("Configuration file change detected, reloading...")
			var newCfg Config
			if err := v.Unmarshal(&newCfg); err != nil {
				log.WithFields(logrus.Fields{
					"file":  e.Name,
					"error": err,
				}).Error("Failed to reload and unmarshal configuration file, continuing with previous settings")
				return
			}

			mu.Lock()
			instance = &newCfg
			mu.Unlock()
			log.WithField("file", e.Name).Info("Configuration successfully reloaded")
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
		log.WithField("error", err).Error("Failed to get configuration instance due to an error during initialization")
		return nil
	}

	mu.RLock()
	defer mu.RUnlock()
	return instance
}

func setViperDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.run_mode", "debug")
	v.SetDefault("server.shutdown_timeout", "30s")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.logger.level", "info")

	v.SetDefault("storage.connect_timeout", "5s")
	v.SetDefault("storage.shutdown_timeout", "15s")
}

func initializeReaderConfig(v *viper.Viper) {
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("internal/config")
	v.AddConfigPath(".")
}
