package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	request_context "main/internal/context"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/mapper"
	"main/internal/response"
	jwtutil "main/internal/security/jwt"
	"main/internal/services"
)

type UserController struct {
	userService services.UserService
	authService services.AuthService
}

func NewUserController(userService services.UserService, authService services.AuthService) *UserController {
	return &UserController{userService: userService, authService: authService}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var input dto.CreateUserInput
	if err := decoder(r, &input); err != nil {
		return err
	}
	err := input.Validate()
	if err != nil {
		return err
	}

	user := mapper.CreateUserInputToModel(&input)
	if user == nil {
		return domain.InvalidRequestError("invalid request body", nil)
	}

	// if err := c.userService.Create(r.Context(), user); err != nil {
	// 	log.Println("Create error: ", err.Error())
	// 	return err
	// }
	_ = c.userService.Create(r.Context(), user)

	log.Println("User: ", user)
	response.WriteJSON(w, r, http.StatusCreated, "user created", mapper.UserModelToDTO(user), nil)
	return nil
}

func (c *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) error {

	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.InvalidRequestError("invalid user context", nil)
	}

	user, err := c.userService.GetByID(r.Context(), userID)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusOK, "user found", mapper.UserModelToDTO(user), nil)
	return nil
}

func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := c.userService.List(r.Context())
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusOK, "users found", mapper.UserModelsToDTO(users), nil)
	return nil
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) error {
	var input dto.LoginInput
	if err := decoder(r, &input); err != nil {
		return err
	}

	user, accessToken, refreshToken, err := c.authService.Authenticate(r.Context(), input.Username, input.Password)
	if err != nil {
		return err
	}

	jwtutil.JwtUtil.SetAuthCookie(w, accessToken, refreshToken)
	response.WriteJSON(w, r, http.StatusOK, "login successful", dto.AuthResult{
		User: *mapper.UserModelToDTO(user),
	}, nil)
	return nil
}

func (c *UserController) RefreshAccessToken(w http.ResponseWriter, r *http.Request) error {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return nil
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		return domain.UnauthorizedError("Refresh token required")
	}
	refreshToken := cookie.Value

	accessToken, err := c.authService.RefreshAccessToken(r.Context(), id, refreshToken)
	if err != nil {
		return err
	}

	jwtutil.JwtUtil.SetAuthCookie(w, accessToken, refreshToken)
	writeJSON(w, http.StatusOK, map[string]interface{}{"message": "access token refreshed", "status": true})
	return nil
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.InvalidRequestError("invalid user context", nil)
	}

	var input dto.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return err
	}
	input.UserID = userID

	user := mapper.UpdateUserInputToModel(&input)
	if user == nil {
		return domain.InvalidRequestError("invalid request body", nil)
	}

	if err := c.userService.UpdateDetails(r.Context(), user); err != nil {
		return err
	}

	updated, err := c.userService.GetByID(r.Context(), userID)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, mapper.UserModelToDTO(updated))
	return nil
}

type updatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (c *UserController) UpdatePassword(w http.ResponseWriter, r *http.Request) error {
	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.InvalidRequestError("invalid user context", nil)
	}

	var input updatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return domain.InvalidRequestError("invalid request body", err)
	}

	if err := c.authService.UpdatePassword(r.Context(), userID, input.OldPassword, input.NewPassword); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	return nil
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.InvalidRequestError("invalid user context", nil)
	}

	if err := c.userService.Delete(r.Context(), userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func parseUUIDParam(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	idStr := chi.URLParam(r, key)
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return uuid.Nil, false
	}
	return id, true
}

func decoder(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return domain.InvalidRequestError("invalid request body", err)
	}
	return nil
}
