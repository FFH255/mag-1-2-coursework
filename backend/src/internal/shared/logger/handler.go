package logger

import (
	"context"
	"io"
	"log/slog"
)

type handler struct {
	next slog.Handler
}

func newHandler(w io.Writer, level slog.Level) *handler {
	next := slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: removeMsgKey,
	})

	return &handler{next: next}
}

func (h *handler) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *handler) Handle(ctx context.Context, rec slog.Record) error {
	fields := getFieldsFromContext(ctx)
	for k, v := range fields {
		rec.Add(k, v)
	}

	return h.next.Handle(ctx, rec)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{next: h.next.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{next: h.next.WithGroup(name)}
}

func removeMsgKey(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		return slog.Attr{}
	}

	return a
}
