package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	LogLevel string
	Key      string

	logger struct {
		zapLogger *zap.Logger
	}

	config interface {
		LogLevel() LogLevel
		AsJSON() bool
	}
)

const (
	LogLevelDebug     LogLevel = "DEBUG"
	LogLevelInfo      LogLevel = "INFO"
	LogLevelWarning   LogLevel = "WARNING"
	LogLevelError     LogLevel = "ERROR"
	LogLevelEmergency LogLevel = "EMERGENCY"
)

const (
	KeyTraceID Key = "trace_id"
	KeyUserID  Key = "user_id"
)

var (
	initOnce     sync.Once
	globalLogger *logger
)

// Init -
func Init(cfg config) {
	initOnce.Do(func() {
		dynamicLevel := zap.NewAtomicLevelAt(parseLevel(string(cfg.LogLevel())))
		encoderCfg := buildProductionEncoderConfig()

		var encoder zapcore.Encoder
		if cfg.AsJSON() {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			dynamicLevel,
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})
}

func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "emergency":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}
