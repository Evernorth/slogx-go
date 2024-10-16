package slogx

import (
	"context"
	"log/slog"
)

type contextAttrsKey struct{}

// ContextWithAttrs adds one or more slog.Attr objects to the provided Context.  A new Context containing the new
// Attrs is returned.
func ContextWithAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	// Get the Attrs map from the Context
	attrMap := getAttrMap(ctx)

	// Copy the map so we do not modify the original
	newAttrMap := make(map[string]slog.Attr)
	for key, value := range *attrMap {
		newAttrMap[key] = value
	}

	// Add the new Attrs to the copy of the map
	for _, newAttr := range newAttrs {
		newAttrMap[newAttr.Key] = newAttr
	}

	// Store the copy in the new Context and return it
	return context.WithValue(ctx, contextAttrsKey{}, &newAttrMap)
}

// getAttrMap returns the slog.Attrs map from the provided Context.
func getAttrMap(ctx context.Context) *map[string]slog.Attr {

	// Read the map from the Context if it exists and return it
	value := ctx.Value(contextAttrsKey{})
	if value != nil {
		attrMap, ok := value.(*map[string]slog.Attr)
		if !ok {
			panic("Could not cast context attrs to *map[string]slog.Attr")
		}
		return attrMap

	} else {

		// Create a new map and return it if the map does not exist in the Context
		attrMap := make(map[string]slog.Attr)
		return &attrMap
	}
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

	// Convert the map to a slice of Attrs
	attrs := make([]slog.Attr, 0, len(attrMap))
	for _, value := range attrMap {
		attrs = append(attrs, value)
	}

	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
