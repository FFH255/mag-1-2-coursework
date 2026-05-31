package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

const (
	errorCtxKey       = "error"
	requestIDCtxKey   = "request_id"
	userIDCtxKey      = "user_id"
	handlerNameCtxKey = "handler_name"
	statusCodeCtxKey  = "status_code"
	msgCtxKey         = "message"
)

type Logger struct {
	logger *slog.Logger
	w      io.Writer
}

func (l Logger) Debug(ctx context.Context, args ...any) {
	l.logger.DebugContext(ctx, "", args...)
}

func (l Logger) Info(ctx context.Context, args ...any) {
	l.logger.InfoContext(ctx, "", args...)
}

func (l Logger) Warning(ctx context.Context, args ...any) {
	l.logger.WarnContext(ctx, "", args...)
}

func (l Logger) Error(ctx context.Context, args ...any) {
	l.logger.ErrorContext(ctx, "", args...)
}

func (l Logger) WithField(ctx context.Context, k string, v any) context.Context {
	return l.WithFields(ctx, map[string]any{k: v})
}

func (l Logger) WithFields(ctx context.Context, fields map[string]any) context.Context {
	existing := getFieldsFromContext(ctx)

	merged := make(map[string]any, len(existing)+len(fields))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range fields {
		merged[k] = v
	}

	return contextWithFields(ctx, merged)
}

func (l Logger) WithError(ctx context.Context, err error) context.Context {
	return l.WithField(ctx, errorCtxKey, err.Error())
}

func (l Logger) WithRequestID(ctx context.Context, requestID string) context.Context {
	return l.WithField(ctx, requestIDCtxKey, requestID)
}

func (l Logger) WithUserID(ctx context.Context, userID int64) context.Context {
	return l.WithField(ctx, userIDCtxKey, userID)
}

func (l Logger) WithHandlerName(ctx context.Context, name string) context.Context {
	return l.WithField(ctx, handlerNameCtxKey, name)
}

func (l Logger) WithStatusCode(ctx context.Context, statusCode int) context.Context {
	return l.WithField(ctx, statusCodeCtxKey, statusCode)
}

func (l Logger) WithMsg(ctx context.Context, msg string) context.Context {
	return l.WithField(ctx, msgCtxKey, msg)
}

func New(useFile bool, path string) *Logger {
	w := os.Stdout

	if useFile {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		w = f
	}

	return &Logger{
		logger: slog.New(newHandler(w, slog.LevelDebug)),
	}
}
