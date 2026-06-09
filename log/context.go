package log

import (
	"context"
	"log/slog"
)

// ctxKey is the private context key used to store request-scoped slog attributes.
type ctxKey struct{}

// ContextWith stores additional slog attributes in ctx, accumulating any that are
// already present. The stored attributes are injected into every log record produced
// by a logger created with New when the context is passed to a *ContextXxx method.
func ContextWith(ctx context.Context, attrs ...slog.Attr) context.Context {
	if ctx == nil {
		return ctx
	}
	existing := attrsFromContext(ctx)
	merged := make([]slog.Attr, len(existing), len(existing)+len(attrs))
	copy(merged, existing)
	merged = append(merged, attrs...)
	return context.WithValue(ctx, ctxKey{}, merged)
}

func attrsFromContext(ctx context.Context) []slog.Attr {
	v, _ := ctx.Value(ctxKey{}).([]slog.Attr)
	return v
}

// contextHandler wraps a slog.Handler and enriches each record with the
// request-scoped attributes stored in the context via ContextWith.
type contextHandler struct {
	slog.Handler
}

// Handle injects context-scoped attributes into r before delegating to the base handler.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs := attrsFromContext(ctx); len(attrs) > 0 {
		r.AddAttrs(attrs...)
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new contextHandler so that context injection is preserved
// when the caller uses logger.With(...).
func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{h.Handler.WithAttrs(attrs)}
}

// WithGroup returns a new contextHandler so that context injection is preserved
// when the caller uses logger.WithGroup(...).
func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{h.Handler.WithGroup(name)}
}
