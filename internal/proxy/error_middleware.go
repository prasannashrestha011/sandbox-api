package proxy

import (
	"log"
	"main/internal/domain"
	"main/internal/response"
	"net/http"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := response.NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		err := rw.Error()
		if err != nil {
			log.Println("Error:", err.Error())
			if !rw.WroteHeader() {
				appErr, ok := err.(*domain.AppError)
				if !ok {
					appErr = domain.InternalError(err)
				}
				response.WriteJSONError(
					w,
					r,
					appErr,
				)
			}
		}

	})
}
