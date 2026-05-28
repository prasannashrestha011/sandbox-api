package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
	"main/internal/proxy"
)

// RegisterUserRoutes wires user endpoints into the router.
func RegisterUserRoutes(r chi.Router, controller *controllers.UserController) {
	r.Route("/users", func(ur chi.Router) {
		ur.Use(proxy.AuthMiddleware)
		ur.Get("/", controller.ListUsers)
		ur.Get("/{id}", controller.GetUserByID)
		ur.Patch("/{id}", controller.UpdateUser)
		ur.Delete("/{id}", controller.DeleteUser)
	})
	r.Route("/auth", func(ar chi.Router) {
		ar.Post("/login", controller.Login)
		ar.Post("/", controller.CreateUser)
		ar.Post("/{id}/refresh", controller.RefreshAccessToken)
		ar.Patch("/{id}/password", controller.UpdatePassword)
	})
}
