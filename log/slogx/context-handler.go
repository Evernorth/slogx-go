package slogx

import (
	"context"
	"log/slog"
)

type contextAttrsKey struct{}

// ContextWithAttrs adds one or more slog.Attr objects to the provided Context.  A new Context containing the new
// Attrs is returned.  The behavior of this function is modeled after the context.WithValue
// function.
func ContextWithAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	attrs := getAttrs(ctx)
	attrs = append(attrs, newAttrs...)

	return context.WithValue(ctx, contextAttrsKey{}, attrs)
}

// getAttrs returns the contextlogger.Attrs from the provided Context.
func getAttrs(ctx context.Context) []slog.Attr {
	// Create the slice
	var attrs []slog.Attr

	// Read the slice from the Gin Context
	value := ctx.Value(contextAttrsKey{})
	if value != nil {
		attrs, ok := value.([]slog.Attr)
		if !ok {
			panic("Could not cast context attrs to []slog.Attr")
		}
		return attrs
	}

	return attrs
}

// ContextHandler is a slog.Handler that adds slog.Attr objects from the provided Context to the slog.Record.
type ContextHandler struct {
	slog.Handler
}

// NewContextHandler returns a new ContextHandler that wraps the provided slog.Handler.
func NewContextHandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{
		Handler: handler,
	}
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := getAttrs(ctx)
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
