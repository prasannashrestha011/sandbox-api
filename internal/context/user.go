package request_context

import (
	"context"
	"main/internal/types"

	"github.com/google/uuid"
)

// setter

// Sets user id in request context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func WithRole(ctx context.Context, role types.Role) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
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

func UserRole(ctx context.Context) (types.Role, bool) {
	role, ok := ctx.Value(userRoleKey).(types.Role)
	return role, ok
}
