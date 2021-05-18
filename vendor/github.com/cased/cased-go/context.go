package cased

import (
	"context"
)

type contextKey int

// ContextKey ...
var ContextKey = contextKey(0)

// GetContextFromContext ...
func GetContextFromContext(ctx context.Context) AuditEvent {
	if context, ok := ctx.Value(ContextKey).(AuditEvent); ok {
		return context
	}
	return nil
}
