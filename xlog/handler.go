package xlog

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"log/slog"
	"os"
	"slices"
)

type handler struct {
	slog.Handler
	keys []xctx.ContextKey
}

func NewPreConfiguredHandler(keys ...xctx.ContextKey) slog.Handler {
	combinedKeys := append(keys, xctx.Transaction)

	slices.Sort(combinedKeys)
	combinedKeys = slices.Compact(combinedKeys)

	return NewCustomHandler(combinedKeys...)
}

func NewCustomHandler(keys ...xctx.ContextKey) slog.Handler {
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
