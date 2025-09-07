package sloglog

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

// ToContext adds the logger to the context.
// If a logger is already present, it is replaced.
func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext extracts the logger from the context.
// If no logger is found, a new one is created, using slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}) as the base.
func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(ctxKey{}).(*slog.Logger)
	if !ok {
		return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
		}))
	}
	return logger
}
