package request_context

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// setter

// Sets user id in request context
func WithUserID(ctx context.Context, userID string) context.Context {
	log.Println("Setting user context", userID)
	return context.WithValue(ctx, userIDKey, userID)
}

// getter

// Retrieves user id from request context
func UserID(ctx context.Context) (uuid.UUID, bool) {
	raw, ok := ctx.Value(userIDKey).(string)
	if !ok || raw == "" {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false
	}
	return id, ok
}
