package logger

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal"
)

type StabLogger struct{}

func (s StabLogger) Debug(ctx context.Context, args ...any) {}

func (s StabLogger) Info(ctx context.Context, args ...any) {}

func (s StabLogger) Warning(ctx context.Context, args ...any) {}

func (s StabLogger) Error(ctx context.Context, args ...any) {}

func (s StabLogger) WithField(ctx context.Context, k string, v any) context.Context {
	return ctx
}

func (s StabLogger) WithFields(ctx context.Context, fields map[string]any) context.Context {
	return ctx
}

func (s StabLogger) WithError(ctx context.Context, err error) context.Context {
	return ctx
}

func (s StabLogger) WithRequestID(ctx context.Context, requestID string) context.Context {
	return ctx
}

func (s StabLogger) WithUserID(ctx context.Context, userID int64) context.Context {
	return ctx
}

func (s StabLogger) WithHandlerName(ctx context.Context, name string) context.Context {
	return ctx
}

func (s StabLogger) WithStatusCode(ctx context.Context, statusCode int) context.Context {
	return ctx
}

func (s StabLogger) WithMsg(ctx context.Context, msg string) context.Context {
	return ctx
}

func NewStab() internal.Logger {
	return StabLogger{}
}
