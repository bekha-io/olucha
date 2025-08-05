package config

import (
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	once sync.Once
	c    *Config
)

// inject secrets from environment variables to the config struct
func injectSecrets() {
	c.DB.DSN = viper.GetString("db_dsn")
}

func Get() *Config {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		viper.WatchConfig()
		viper.ReadInConfig()

		if err := viper.Unmarshal(c); err != nil {
			logger.Fatal().Err(err).Msg("failed to unmarshal config")
		}

		viper.OnConfigChange(func(e fsnotify.Event) {
			if err := viper.Unmarshal(c); err != nil {
				logger.Error().Err(err).Msg("failed to unmarshal config")
			}
			injectSecrets()
		})
	})

	return c
}

type Config struct {
	DB DB
}

type DB struct {
	DSN string
}
