package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
	"main/internal/proxy"
	"main/internal/response"
)

// RegisterSandboxRoutes wires sandbox endpoints into the router.
func RegisterSandboxRoutes(r chi.Router, controller *controllers.SandboxController) {
	r.Route("/sandbox-template", func(sr chi.Router) {
		sr.Use(proxy.AuthMiddleware)
		response.WrapPost(sr, "/", controller.CreateSandbox)
		response.WrapGet(sr, "/{id}", controller.GetSandboxByID)
		response.WrapGet(sr, "/", controller.ListTemplatesByUser)
		response.WrapPut(sr, "/{id}", controller.UpdateSandbox)
	})
}
