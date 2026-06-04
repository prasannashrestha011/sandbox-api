package pkg

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ExtractParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func ExtractQuery(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
