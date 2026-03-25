package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// Config holds the application configuration values.
type Config struct {
	AppName  string
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
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

// DatabaseConfig contains database connection settings.
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LogConfig controls structured logging behavior.
type LogConfig struct {
	Level string
}

// Load reads configuration from file (optional) and environment variables.
// Environment variables override file values. The APP_ prefix is used for env keys.
func Load(configPath string) (Config, error) {
	if err := loadDotEnvIfPresent(); err != nil {
		return Config{}, err
	}

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
		Database: DatabaseConfig{
			// DB DSN is intentionally env-only to enforce a single source of truth.
			URL:             strings.TrimSpace(os.Getenv("APP_DATABASE_URL")),
			MaxOpenConns:    viper.GetInt("database.maxOpenConns"),
			MaxIdleConns:    viper.GetInt("database.maxIdleConns"),
			ConnMaxLifetime: viper.GetDuration("database.connMaxLifetime"),
		},
		Log: LogConfig{Level: viper.GetString("log.level")},
	}

	normalize(&cfg)
	if cfg.Database.URL == "" {
		return Config{}, fmt.Errorf("missing required environment variable APP_DATABASE_URL")
	}
	return cfg, nil
}

func loadDotEnvIfPresent() error {
	if _, err := os.Stat(".env"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to check .env: %w", err)
	}

	if err := gotenv.Load(".env"); err != nil {
		return fmt.Errorf("failed to load .env: %w", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("app.name", "go-rest-template")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.shutdownTimeout", "10s")
	viper.SetDefault("server.readTimeout", "15s")
	viper.SetDefault("server.writeTimeout", "15s")
	viper.SetDefault("server.idleTimeout", "60s")
	viper.SetDefault("server.allowedOrigins", []string{"http://localhost:3000", "http://127.0.0.1:3000"})
	viper.SetDefault("database.maxOpenConns", 25)
	viper.SetDefault("database.maxIdleConns", 5)
	viper.SetDefault("database.connMaxLifetime", "5m")
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
	if cfg.Database.MaxOpenConns <= 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleConns <= 0 {
		cfg.Database.MaxIdleConns = 5
	}
	if cfg.Database.ConnMaxLifetime <= 0 {
		cfg.Database.ConnMaxLifetime = 5 * time.Minute
	}
}
