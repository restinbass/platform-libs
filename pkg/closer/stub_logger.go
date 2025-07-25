package closer

import (
	"context"

	"go.uber.org/zap"
)

type noopLogger struct{}

func (l *noopLogger) Info(ctx context.Context, msg string, fields ...zap.Field)  {}
func (l *noopLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {}
