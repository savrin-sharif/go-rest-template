package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration values.
type Config struct {
	AppName string
	Server  ServerConfig
	Log     LogConfig
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port            int
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

// LogConfig controls structured logging behavior.
type LogConfig struct {
	Level string
}

// Load reads configuration from file (optional) and environment variables.
// Environment variables override file values. The APP_ prefix is used for env keys.
func Load(configPath string) (Config, error) {
	setDefaults()

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("failed to read config: %w", err)
		}
	}

	cfg := Config{
		AppName: viper.GetString("app.name"),
		Server: ServerConfig{
			Port:            viper.GetInt("server.port"),
			ShutdownTimeout: viper.GetDuration("server.shutdownTimeout"),
			AllowedOrigins:  viper.GetStringSlice("server.allowedOrigins"),
			ReadTimeout:     viper.GetDuration("server.readTimeout"),
			WriteTimeout:    viper.GetDuration("server.writeTimeout"),
			IdleTimeout:     viper.GetDuration("server.idleTimeout"),
		},
		Log: LogConfig{Level: viper.GetString("log.level")},
	}

	normalize(&cfg)
	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("app.name", "go-rest-template")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.shutdownTimeout", "10s")
	viper.SetDefault("server.readTimeout", "15s")
	viper.SetDefault("server.writeTimeout", "15s")
	viper.SetDefault("server.idleTimeout", "60s")
	viper.SetDefault("server.allowedOrigins", []string{"http://localhost:3000", "http://127.0.0.1:3000"})
	viper.SetDefault("log.level", "info")
}

func normalize(cfg *Config) {
	if cfg.Server.ShutdownTimeout <= 0 {
		cfg.Server.ShutdownTimeout = 10 * time.Second
	}
	if cfg.Server.ReadTimeout <= 0 {
		cfg.Server.ReadTimeout = 15 * time.Second
	}
	if cfg.Server.WriteTimeout <= 0 {
		cfg.Server.WriteTimeout = 15 * time.Second
	}
	if cfg.Server.IdleTimeout <= 0 {
		cfg.Server.IdleTimeout = 60 * time.Second
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if len(cfg.Server.AllowedOrigins) == 0 {
		cfg.Server.AllowedOrigins = []string{"*"}
	}
}
