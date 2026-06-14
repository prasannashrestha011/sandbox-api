package controllers

import (
	"encoding/json"
	"log"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/response"
	"main/internal/services"
	"net/http"
)

type SandboxSessionHandler struct {
	sandboxSessionService services.SandboxSessionService
}

func NewSandboxSessionHandler(svc services.SandboxSessionService) *SandboxSessionHandler {
	return &SandboxSessionHandler{
		sandboxSessionService: svc,
	}
}

func (h *SandboxSessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateSessionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid templateID", nil)
	}
	log.Printf("Received request to create session with template ID: %s", req.TemplateID)
	resp, err := h.sandboxSessionService.CreateSession(r.Context(), req.TemplateID)
	if err != nil {
		return err
	}
	log.Printf("Session created successfully with ID: %s", resp.SessionID)
	response.WriteJSON(w, r, http.StatusOK, "create new session", resp, nil)
	return nil

}
