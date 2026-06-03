package proxy

import (
	"log"
	"net/http"

	request_context "main/internal/context"
	"main/internal/enums"
)

// stands before lab routes. User with instructor privileges can only access the lab routes.
func UserTypeMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, ok := request_context.UserType(r.Context())
		log.Println("User type from context: ", userType, "ok: ", ok)
		if !ok {
			http.Error(w, "User type not found in context", http.StatusUnauthorized)
			return
		}
		if userType != enums.UserTypeInstructor {
			http.Error(w, "Forbidden: only instructor can access this resource", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
