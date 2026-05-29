package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"main/internal/mapper"
	"main/internal/services"
	sandbox_type "main/internal/types"
)

type SandboxController struct {
	service services.SandboxService
}

func NewSandboxController(service services.SandboxService) *SandboxController {
	return &SandboxController{service: service}
}

func (c *SandboxController) CreateSandbox(w http.ResponseWriter, r *http.Request) {
	var req sandbox_type.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	sandbox := mapper.SandboxCreateRequestToModel(req, time.Now())

	if err := c.service.Create(r.Context(), sandbox); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create sandbox"})
		return
	}

	writeJSON(w, http.StatusCreated, mapper.SandboxModelToCreateResponse(sandbox))
}

func (c *SandboxController) GetSandboxByID(w http.ResponseWriter, r *http.Request) {
	idStr := extractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid sandbox id"})
		return
	}

	sandbox, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "sandbox not found"})
		return
	}

	writeJSON(w, http.StatusOK, sandbox)
}

func (c *SandboxController) GetSandboxBySessionID(w http.ResponseWriter, r *http.Request) {
	sessionStr := extractParam(r, "sessionId")
	sessionID, err := uuid.Parse(sessionStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid session id"})
		return
	}

	sandbox, err := c.service.GetBySessionID(r.Context(), sessionID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "sandbox not found"})
		return
	}

	writeJSON(w, http.StatusOK, sandbox)
}

func (c *SandboxController) ListSandboxesByUser(w http.ResponseWriter, r *http.Request) {
	userID := extractQuery(r, "user_id")
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id is required"})
		return
	}

	items, err := c.service.ListByUserID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list sandboxes"})
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func (c *SandboxController) UpdateSandboxStatus(w http.ResponseWriter, r *http.Request) {
	idStr := extractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid sandbox id"})
		return
	}

	var req sandbox_type.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Status == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "status is required"})
		return
	}

	if err := c.service.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update status"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (c *SandboxController) DeleteSandbox(w http.ResponseWriter, r *http.Request) {
	idStr := extractParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid sandbox id"})
		return
	}

	if err := c.service.Delete(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete sandbox"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
