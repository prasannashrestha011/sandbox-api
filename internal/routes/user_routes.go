package routes

import (
	"github.com/go-chi/chi/v5"

	"main/internal/controllers"
	"main/internal/response"
)

// RegisterUserRoutes wires user endpoints into the router.
func RegisterUserRoutes(r chi.Router, controller *controllers.UserController) {
	r.Route("/users", func(ur chi.Router) {
		response.WrapGet(ur, "/", controller.ListUsers)
		response.WrapGet(ur, "/{id}", controller.GetUserByID)
		response.WrapPatch(ur, "/{id}", controller.UpdateUser)
		response.WrapDelete(ur, "/{id}", controller.DeleteUser)
	})
	r.Route("/auth", func(ar chi.Router) {
		response.WrapPost(ar, "/register", controller.CreateUser)
		response.WrapPost(ar, "/login", controller.Login)
		response.WrapPost(ar, "/", controller.CreateUser)
		response.WrapPost(ar, "/{id}/refresh", controller.RefreshAccessToken)
		response.WrapPatch(ar, "/{id}/password", controller.UpdatePassword)
	})
}
