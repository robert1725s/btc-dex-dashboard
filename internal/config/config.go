package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Job      JobConfig      `mapstructure:"job"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type JobConfig struct {
	IntervalSeconds int `mapstructure:"interval_seconds"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// デフォルト値
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.path", "dev.db")
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:5173"})
	viper.SetDefault("job.interval_seconds", 2)

	// 環境変数での上書きを許可
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// config.yaml がなくてもデフォルト値で動作
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
