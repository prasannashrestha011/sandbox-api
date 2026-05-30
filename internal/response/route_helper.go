package response

import (
	"github.com/go-chi/chi/v5"
)

// Helper functions for wrapping handlers with error handling

func WrapGet(r chi.Router, path string, h HandlerFunc) {
	r.Get(path, Wrap(h))
}

func WrapPost(r chi.Router, path string, h HandlerFunc) {
	r.Post(path, Wrap(h))
}

func WrapPut(r chi.Router, path string, h HandlerFunc) {
	r.Put(path, Wrap(h))
}

func WrapPatch(r chi.Router, path string, h HandlerFunc) {
	r.Patch(path, Wrap(h))
}

func WrapDelete(r chi.Router, path string, h HandlerFunc) {
	r.Delete(path, Wrap(h))
}
