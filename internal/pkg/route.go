package pkg

import (
	"main/internal/proxy"

	"github.com/go-chi/chi/v5"
)

func NewWebSocket(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(proxy.AuthMiddleware)
		r.Get("/ws", Manager.HandleConnection)
	})
}
