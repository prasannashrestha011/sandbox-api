package routes

import (
	"main/internal/controllers"
	"main/internal/proxy"
	"main/internal/response"

	"github.com/go-chi/chi/v5"
)

func RegisterDockerImageRoutes(router *chi.Mux, controller *controllers.DockerImageController) {
	router.Route("/docker-images", func(r chi.Router) {
		r.Use(proxy.AuthMiddleware)
		r.Use(proxy.AdminMiddleware)
		response.WrapPost(r, "/", controller.CreateImage)
		response.WrapGet(r, "/", controller.ListImages)
	})
}
