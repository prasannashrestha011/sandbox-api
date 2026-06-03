package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
	"main/internal/proxy"
	"main/internal/response"
)

// RegisterLabRoutes wires lab endpoints into the router.
func RegisterLabRoutes(r chi.Router, controller *controllers.LabController) {
	r.Route("/labs", func(sr chi.Router) {
		sr.Use(proxy.AuthMiddleware)
		sr.Use(proxy.UserTypeMiddlware)
		response.WrapPost(sr, "/", controller.CreateLab)
		response.WrapGet(sr, "/{id}", controller.GetLabByID)
		response.WrapDelete(sr, "/{id}", controller.DeleteLab)
	})
}
