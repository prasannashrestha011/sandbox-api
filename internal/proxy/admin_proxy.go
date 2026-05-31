package proxy

import (
	"net/http"

	request_context "main/internal/context"
	"main/internal/enums"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		role, ok := request_context.UserRole(ctx)
		if !ok {
			http.Error(w, "Missing role", http.StatusUnauthorized)
			return
		}
		if role != enums.RoleAdmin {
			http.Error(w, "Required admin privilege", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
