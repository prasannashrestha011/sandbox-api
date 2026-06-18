package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
	"main/internal/proxy"
	"main/internal/response"
)

// RegisterSandboxRoutes wires sandbox endpoints into the router.
func RegisterSandboxRoutes(r chi.Router, template *controllers.SandboxController, session *controllers.SandboxSessionHandler) {
	r.Route("/sandbox-template", func(sr chi.Router) {
		sr.Use(proxy.AuthMiddleware)
		response.WrapPost(sr, "/", template.CreateSandbox)
		response.WrapGet(sr, "/{id}", template.GetSandboxByID)
		response.WrapGet(sr, "/", template.ListTemplatesByUser)
		response.WrapPut(sr, "/{id}", template.UpdateSandbox)
	})
	r.Route("/sandbox/session", func(r chi.Router) {
		r.Use(proxy.AuthMiddleware)
		response.WrapPost(r, "/", session.CreateSession)
		response.WrapPost(r, "/{id}/execute", session.ExecuteCode)
	})
}
