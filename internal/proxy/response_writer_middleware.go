package proxy

import (
	"main/internal/response"
	"net/http"
)

func ResponseWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &response.ResponseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
	})
}
