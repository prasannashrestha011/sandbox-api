package proxy

import (
	"log"
	"net/http"

	request_context "main/internal/context"
	"main/internal/enums"
)

// stands before lab routes. User with instructor privileges can only access the lab routes.
// hides the lab routes from students
func UserTypeMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType, ok := request_context.UserType(r.Context())
		log.Println("User type from context: ", userType, "ok: ", ok)
		if !ok {
			http.NotFound(w, r)
			return
		}
		if userType != enums.UserTypeInstructor {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
