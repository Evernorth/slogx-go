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
	attrMap := getAttrMap(ctx)
	for _, newAttr := range newAttrs {
		attrMap[newAttr.Key] = newAttr
	}

	return context.WithValue(ctx, contextAttrsKey{}, attrMap)
}

// getAttrMap returns the slog.Attrs from the provided Context.
func getAttrMap(ctx context.Context) map[string]slog.Attr {
	// Create the map
	attrMap := make(map[string]slog.Attr)

	// Read the slice from the Context
	value := ctx.Value(contextAttrsKey{})
	if value != nil {
		attrMap, ok := value.(map[string]slog.Attr)
		if !ok {
			panic("Could not cast context attrs to []slog.Attr")
		}
		return attrMap
	}

	return attrMap
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
	attrMap := getAttrMap(ctx)

	// Convert the map to a slice
	attrs := make([]slog.Attr, 0, len(attrMap))
	for _, value := range attrMap {
		attrs = append(attrs, value)
	}

	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
