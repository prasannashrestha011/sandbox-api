package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
)

// RegisterUserRoutes wires user endpoints into the router.
func RegisterUserRoutes(r chi.Router, controller *controllers.UserController) {
	r.Route("/users", func(ur chi.Router) {
		ur.Post("/", controller.CreateUser)
		ur.Post("/login", controller.Login)
		ur.Get("/", controller.ListUsers)
		ur.Get("/{id}", controller.GetUserByID)
		ur.Patch("/{id}", controller.UpdateUser)
		ur.Patch("/{id}/password", controller.UpdatePassword)
		ur.Delete("/{id}", controller.DeleteUser)
	})
}
