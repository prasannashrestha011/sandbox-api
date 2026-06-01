package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"main/internal/controllers/mapper"
	"main/internal/domain"
	"main/internal/dto"
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

	writeJSON(w, http.StatusCreated, mapper.SandboxServiceModelToCreateResponse(sandbox))
	return nil
}

func (c *SandboxController) GetSandboxByID(w http.ResponseWriter, r *http.Request) error {
	idStr := extractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return domain.InvalidRequestError("invalid sandbox id", nil)
	}

	sandbox, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, sandbox)
	return nil
}

func (c *SandboxController) GetSandboxBySessionID(w http.ResponseWriter, r *http.Request) error {
	sessionStr := extractParam(r, "sessionId")
	sessionID, err := uuid.Parse(sessionStr)
	if err != nil {
		return err
	}

	sandbox, err := c.service.GetBySessionID(r.Context(), sessionID)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, sandbox)
	return nil
}

func (c *SandboxController) ListSandboxesByUser(w http.ResponseWriter, r *http.Request) error {
	userID := extractQuery(r, "user_id")
	if userID == "" {
		return domain.InvalidRequestError("user_id is required", nil)
	}

	items, err := c.service.ListByUserID(r.Context(), userID)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, items)
	return nil
}

func (c *SandboxController) UpdateSandboxStatus(w http.ResponseWriter, r *http.Request) error {
	idStr := extractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return domain.InvalidRequestError("invalid sandbox id", nil)
	}

	var req dto.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", nil)
	}
	if req.Status == "" {
		return domain.InvalidRequestError("status is required", nil)
	}

	if err := c.service.UpdateStatus(r.Context(), id, req.Status); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	return nil
}

func (c *SandboxController) DeleteSandbox(w http.ResponseWriter, r *http.Request) error {
	idStr := extractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return domain.InvalidRequestError("invalid sandbox id", nil)
	}

	if err := c.service.Delete(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *SandboxController) ExecuteCode(w http.ResponseWriter, r *http.Request) error {
	containerIDStr := extractParam(r, "id")

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

	writeJSON(w, http.StatusOK, map[string]string{"result": result})
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func extractParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func extractQuery(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
