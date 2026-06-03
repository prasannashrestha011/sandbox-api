package proxy

import (
	"log"
	"net/http"
	"time"

	request_context "main/internal/context"
	"main/internal/enums"
	jwtutil "main/internal/security/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			log.Println("Missing auth cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claim, err := jwtutil.JwtUtil.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}

		if claim.ExpiresAt != nil && claim.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}
		ctx = request_context.WithUserID(ctx, claim.UserID)
		ctx = request_context.WithRole(ctx, enums.Role(claim.Role))
		ctx = request_context.WithUserType(ctx, claim.UserType)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
