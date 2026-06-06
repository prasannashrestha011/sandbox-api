package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	request_context "main/internal/context"
	"main/internal/controllers/mapper"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/pkg"
	"main/internal/response"
	"main/internal/services"
)

type SandboxController struct {
	service services.SandboxService
}

func NewSandboxController(service services.SandboxService) *SandboxController {
	return &SandboxController{service: service}
}

func (c *SandboxController) CreateSandbox(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	sandbox, err := mapper.SandboxCreateRequestToServiceModel(req, r.Context(), time.Now())
	if err != nil {
		return err
	}

	if err := c.service.Create(r.Context(), req.ImageID, sandbox); err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusCreated, "sandbox created", mapper.SandboxServiceModelToCreateResponse(sandbox), nil)
	return nil
}

func (c *SandboxController) GetSandboxByID(w http.ResponseWriter, r *http.Request) error {
	idStr := pkg.ExtractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return domain.InvalidRequestError("invalid sandbox id", nil)
	}

	sandbox, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusOK, "sandbox found", mapper.SandboxServiceModelToCreateResponse(sandbox), nil)
	return nil
}

func (c *SandboxController) GetSandboxBySessionID(w http.ResponseWriter, r *http.Request) error {
	sessionStr := pkg.ExtractParam(r, "sessionId")
	sessionID, err := uuid.Parse(sessionStr)
	if err != nil {
		return err
	}

	sandbox, err := c.service.GetBySessionID(r.Context(), sessionID)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusOK, "sandbox found", mapper.SandboxServiceModelToCreateResponse(sandbox), nil)
	return nil
}

func (c *SandboxController) ListSandboxesByUser(w http.ResponseWriter, r *http.Request) error {
	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.UnauthorizedError("user id not found in context")
	}

	items, err := c.service.ListByUserID(r.Context(), userID.String())
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "sandboxes found", items, nil)
	return nil
}

func (c *SandboxController) DeleteSandbox(w http.ResponseWriter, r *http.Request) error {
	idStr := pkg.ExtractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return domain.InvalidRequestError("invalid sandbox id", nil)
	}

	if err := c.service.Delete(r.Context(), id.String()); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *SandboxController) ExecuteCode(w http.ResponseWriter, r *http.Request) error {
	containerIDStr := pkg.ExtractParam(r, "id")

	var req dto.ExecuteCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", nil)
	}
	if req.Code == "" {
		return domain.InvalidRequestError("code is required", nil)
	}

	result, err := c.service.ExecuteCode(r.Context(), containerIDStr, &req)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusOK, "code executed", map[string]string{"result": result}, nil)
	return nil
}
