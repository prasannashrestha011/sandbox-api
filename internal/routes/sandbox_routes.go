package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
	"main/internal/proxy"
	"main/internal/response"
)

// RegisterSandboxRoutes wires sandbox endpoints into the router.
func RegisterSandboxRoutes(r chi.Router, controller *controllers.SandboxController) {
	r.Route("/sandboxes", func(sr chi.Router) {
		sr.Use(proxy.AuthMiddleware)
		response.WrapPost(sr, "/", controller.CreateSandbox)
		response.WrapGet(sr, "/{id}", controller.GetSandboxByID)
		response.WrapGet(sr, "/", controller.ListSandboxesByUser)
		response.WrapGet(sr, "/session/{sessionId}", controller.GetSandboxBySessionID)
		response.WrapPost(sr, "/{id}/execute", controller.ExecuteCode)
		response.WrapPatch(sr, "/{id}/status", controller.UpdateSandboxStatus)
		response.WrapDelete(sr, "/{id}", controller.DeleteSandbox)
	})
}
