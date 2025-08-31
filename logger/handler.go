package logger

import (
	"context"
	"log/slog"
	"os"
	"slices"
)

type ContextKey string

const (
	Transaction ContextKey = "tx_id"
)

type handler struct {
	slog.Handler
	keys []ContextKey
}

func NewPreConfiguredHandler(keys ...ContextKey) slog.Handler {
	combinedKeys := append(keys, Transaction)

	slices.Sort(combinedKeys)
	combinedKeys = slices.Compact(combinedKeys)

	return NewCustomHandler(combinedKeys...)
}

func NewCustomHandler(keys ...ContextKey) slog.Handler {
	return &handler{
		Handler: slog.NewJSONHandler(os.Stdout, nil),
		keys:    keys,
	}
}

func (ch *handler) Handle(ctx context.Context, r slog.Record) error {
	for _, key := range ch.keys {
		rawValue := ctx.Value(key)
		if rawValue == nil {
			continue
		}

		switch value := rawValue.(type) {
		case string:
			r.AddAttrs(slog.String(string(key), value))
		case int:
			r.AddAttrs(slog.Int64(string(key), int64(value)))
		case int64:
			r.AddAttrs(slog.Int64(string(key), value))
		}
	}

	return ch.Handler.Handle(ctx, r)
}
