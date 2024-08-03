package utils

import (
	"context"

	"go.uber.org/zap"
)

func getZapLoggerLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func InitializeLogger(logLevel string) (*zap.Logger, func(), error) {
	zapLoggerConfig := zap.NewProductionConfig()

	// Default set details debug level | set level from configs
	// Filter by log level
	zapLoggerConfig.Level = getZapLoggerLevel(logLevel)

	logger, err := zapLoggerConfig.Build()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		// deliberately ignore the returned error here
		_ = logger.Sync()
	}

	return logger, cleanup, err
}

func LoggerWithContext(_ context.Context, logger *zap.Logger) *zap.Logger {
	//TODO: Add request ID to context
	return logger
}
