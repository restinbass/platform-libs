package logger

import (
	"context"

	"go.uber.org/zap"
)

// Logger -
func Logger() *logger {
	return globalLogger
}

// Sync -
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}

	return nil
}

// With -
func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext -
func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(KeyTraceID).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(KeyTraceID), traceID))
	}

	if userID, ok := ctx.Value(KeyUserID).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(KeyUserID), userID))
	}

	return fields
}
