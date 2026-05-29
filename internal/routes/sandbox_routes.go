package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
	"main/internal/proxy"
)

// RegisterSandboxRoutes wires sandbox endpoints into the router.
func RegisterSandboxRoutes(r chi.Router, controller *controllers.SandboxController) {
	r.Route("/sandboxes", func(sr chi.Router) {
		sr.Use(proxy.AuthMiddleware)
		sr.Post("/", controller.CreateSandbox)
		sr.Get("/", controller.ListSandboxesByUser)
		sr.Get("/{id}", controller.GetSandboxByID)
		sr.Get("/session/{sessionId}", controller.GetSandboxBySessionID)
		sr.Patch("/{id}/status", controller.UpdateSandboxStatus)
		sr.Delete("/{id}", controller.DeleteSandbox)
	})
}
