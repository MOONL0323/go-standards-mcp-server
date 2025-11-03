package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"go-standards-mcp-server/internal/analyzer"
	"go-standards-mcp-server/internal/config"
	"go-standards-mcp-server/internal/mcp"
	"go-standards-mcp-server/internal/usercontext"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	configPath = flag.String("config", "", "Path to config file")
	mode       = flag.String("mode", "", "Server mode: stdio or http")
	port       = flag.Int("port", 0, "HTTP server port")
	logLevel   = flag.String("log-level", "", "Log level: debug, info, warn, error")
	version    = flag.Bool("version", false, "Print version and exit")
)

const (
	appName    = "go-standards-mcp-server"
	appVersion = "0.5.0"
)

func main() {
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Override config with command-line flags
	if *mode != "" {
		cfg.Server.Mode = *mode
	}
	if *port != 0 {
		cfg.Server.Port = *port
	}
	if *logLevel != "" {
		cfg.Log.Level = *logLevel
	}

	// Initialize logger
	logger, err := initLogger(cfg.Log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting server",
		zap.String("version", appVersion),
		zap.String("mode", cfg.Server.Mode))

	// Initialize session manager for multi-user support
	sessionTimeout := 30 * time.Minute // Default 30 minutes
	if cfg.Server.SessionTimeout > 0 {
		sessionTimeout = time.Duration(cfg.Server.SessionTimeout) * time.Minute
	}
	sessionManager := usercontext.NewSessionManager("./storage", sessionTimeout)
	logger.Info("Session manager initialized",
		zap.Duration("timeout", sessionTimeout))

	// Initialize analyzer
	analyzer, err := analyzer.NewAnalyzer(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize analyzer", zap.Error(err))
	}

	// Initialize MCP server
	server, err := mcp.NewServer(cfg, logger, analyzer, sessionManager)
	if err != nil {
		logger.Fatal("Failed to initialize MCP server", zap.Error(err))
	}

	// Log session statistics
	logger.Info("Server ready",
		zap.Int("active_sessions", sessionManager.GetActiveSessionCount()))

	// Start server
	if err := server.Serve(); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

// initLogger initializes the logger
func initLogger(cfg config.LogConfig) (*zap.Logger, error) {
	// Parse log level
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	// Create config
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          cfg.Format,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{cfg.Output},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Build logger
	return zapConfig.Build()
}
