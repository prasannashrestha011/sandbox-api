package response

import "net/http"

// HandlerFunc is like http.HandlerFunc but returns an error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Wrap converts HandlerFunc to standard http.Handler
func Wrap(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw, ok := w.(*ResponseWriter) // type assert once
		if !ok {
			panic("response writer is not *ResponseWriter")
		}

		if err := h(rw, r); err != nil {
			rw.SetError(err) // automatically store error
		}
	}
}
