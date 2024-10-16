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

	// Get the Attrs map from the Context
	attrMap := getAttrMap(ctx)

	// Copy the map
	newAttrMap := make(map[string]slog.Attr)
	for key, value := range *attrMap {
		newAttrMap[key] = value
	}

	// Add the new Attrs to the copy
	for _, newAttr := range newAttrs {
		newAttrMap[newAttr.Key] = newAttr
	}

	// Store the copy in the new Context
	return context.WithValue(ctx, contextAttrsKey{}, &newAttrMap)
}

// getAttrMap returns the slog.Attrs map from the provided Context.
func getAttrMap(ctx context.Context) *map[string]slog.Attr {

	// Read the map from the Context
	value := ctx.Value(contextAttrsKey{})
	if value != nil {
		attrMap, ok := value.(*map[string]slog.Attr)
		if !ok {
			panic("Could not cast context attrs to []slog.Attr")
		}
		return attrMap
	}

	// Create the map
	attrMap := make(map[string]slog.Attr)
	return &attrMap
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
	attrMap := *getAttrMap(ctx)

	// Convert the map to a slice
	attrs := make([]slog.Attr, 0, len(attrMap))
	for _, value := range attrMap {
		attrs = append(attrs, value)
	}

	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
