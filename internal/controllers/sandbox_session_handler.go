package controllers

import (
	"encoding/json"
	"log"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/pkg"
	"main/internal/response"
	"main/internal/services"
	"net/http"

	"github.com/google/uuid"
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

func (h *SandboxSessionHandler) ExecuteCode(w http.ResponseWriter, r *http.Request) error {
	sessionID := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(sessionID); err != nil {
		return domain.InvalidRequestError("invalid session ID", nil)
	}
	var req dto.SandboxExecReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", nil)
	}
	resp, err := h.sandboxSessionService.ExecuteCommand(r.Context(), sessionID, &req)
	if err != nil {
		return err
	}
	log.Printf("Code executed successfully in session ID: %s", sessionID)
	response.WriteJSON(w, r, http.StatusOK, "code execution result", resp, nil)
	return nil
}
