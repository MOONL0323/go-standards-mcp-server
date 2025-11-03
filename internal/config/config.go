package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
	Analyzer AnalyzerConfig `mapstructure:"analyzer"`
	Linters  LintersConfig  `mapstructure:"linters"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Report   ReportConfig   `mapstructure:"report"`
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Mode           string `mapstructure:"mode"`            // stdio or http
	Port           int    `mapstructure:"port"`
	Host           string `mapstructure:"host"`
	SessionTimeout int    `mapstructure:"session_timeout"` // Session timeout in minutes (default: 30)
}

// LogConfig contains logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Output string `mapstructure:"output"` // stdout, file path
	Format string `mapstructure:"format"` // json, text
}

// AnalyzerConfig contains analyzer configuration
type AnalyzerConfig struct {
	Timeout         time.Duration `mapstructure:"timeout"`
	ConcurrentLimit int           `mapstructure:"concurrent_limit"`
	TempDir         string        `mapstructure:"temp_dir"`
}

// LintersConfig contains linter configurations
type LintersConfig struct {
	GolangciLint GolangciLintConfig `mapstructure:"golangci_lint"`
	Staticcheck  LinterConfig       `mapstructure:"staticcheck"`
	Gosec        LinterConfig       `mapstructure:"gosec"`
	Govet        LinterConfig       `mapstructure:"govet"`
}

// GolangciLintConfig contains golangci-lint specific configuration
type GolangciLintConfig struct {
	Enabled    bool          `mapstructure:"enabled"`
	Timeout    time.Duration `mapstructure:"timeout"`
	ConfigPath string        `mapstructure:"config_path"`
}

// LinterConfig contains generic linter configuration
type LinterConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// StorageConfig contains storage configuration
type StorageConfig struct {
	Type     string         `mapstructure:"type"` // sqlite or postgres
	SQLite   SQLiteConfig   `mapstructure:"sqlite"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

// SQLiteConfig contains SQLite configuration
type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

// PostgresConfig contains PostgreSQL configuration
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

// CacheConfig contains cache configuration
type CacheConfig struct {
	Enabled bool        `mapstructure:"enabled"`
	Type    string      `mapstructure:"type"` // redis, memory
	Redis   RedisConfig `mapstructure:"redis"`
}

// RedisConfig contains Redis configuration
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// ReportConfig contains report configuration
type ReportConfig struct {
	OutputDir string   `mapstructure:"output_dir"`
	Formats   []string `mapstructure:"formats"`
	KeepDays  int      `mapstructure:"keep_days"`
}

// Load loads configuration from file and environment
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("default")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults
	}

	// Override with environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("MCP")

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.mode", "stdio")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("log.format", "json")

	v.SetDefault("analyzer.timeout", "300s")
	v.SetDefault("analyzer.concurrent_limit", 10)
	v.SetDefault("analyzer.temp_dir", "./tmp")

	v.SetDefault("linters.golangci_lint.enabled", true)
	v.SetDefault("linters.golangci_lint.timeout", "5m")
	v.SetDefault("linters.staticcheck.enabled", true)
	v.SetDefault("linters.gosec.enabled", true)
	v.SetDefault("linters.govet.enabled", true)

	v.SetDefault("storage.type", "sqlite")
	v.SetDefault("storage.sqlite.path", "./data/mcp_server.db")

	v.SetDefault("cache.enabled", false)
	v.SetDefault("cache.type", "redis")

	v.SetDefault("report.output_dir", "./reports")
	v.SetDefault("report.formats", []string{"json", "markdown"})
	v.SetDefault("report.keep_days", 30)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server mode
	if c.Server.Mode != "stdio" && c.Server.Mode != "http" {
		return fmt.Errorf("invalid server mode: %s (must be stdio or http)", c.Server.Mode)
	}

	// Validate log level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Log.Level] {
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}

	// Validate storage type
	if c.Storage.Type != "sqlite" && c.Storage.Type != "postgres" {
		return fmt.Errorf("invalid storage type: %s", c.Storage.Type)
	}

	// Create necessary directories
	dirs := []string{
		c.Analyzer.TempDir,
		c.Report.OutputDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
