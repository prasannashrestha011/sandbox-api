package routes

import (
	"main/internal/controllers"
	"main/internal/proxy"
	"main/internal/response"

	"github.com/go-chi/chi/v5"
)

func RegisterWarmPoolRoutes(router *chi.Mux, handler *controllers.WarmPoolHandler) {
	router.Route("/warmpools", func(r chi.Router) {
		r.Use(proxy.AuthMiddleware)
		r.Use(proxy.AdminMiddleware)
		response.WrapPost(r, "/", handler.CreateWarmPool)
	})
}
