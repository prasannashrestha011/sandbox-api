package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	request_context "main/internal/context"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/pkg"
	"main/internal/response"
	"main/internal/services"
	"main/internal/services/mapper"
)

type SandboxController struct {
	service services.SandboxTemplateService
}

func NewSandboxController(service services.SandboxTemplateService) *SandboxController {
	return &SandboxController{service: service}
}

func (c *SandboxController) CreateSandbox(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateTemplateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	sandbox, err := mapper.ToSandboxTemplate(req, r.Context(), time.Now())
	if err != nil {
		return err
	}

	if err := c.service.Create(r.Context(), req.ImageID, sandbox); err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusCreated, "sandbox created", nil, nil)
	return nil
}

func (c *SandboxController) GetSandboxByID(w http.ResponseWriter, r *http.Request) error {
	id := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		return domain.InvalidRequestError("invalid sandbox id", nil)
	}

	sandbox, err := c.service.GetByID(r.Context(), id)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusOK, "sandbox found", mapper.ToSandboxTemplateResponse(sandbox), nil)
	return nil
}

func (c *SandboxController) ListTemplatesByUser(w http.ResponseWriter, r *http.Request) error {
	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.UnauthorizedError("user id not found in context")
	}

	items, err := c.service.ListByUserID(r.Context(), userID.String())
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "sandboxes found", mapper.ToSandboxTemplateResponseList(items), nil)
	return nil
}

func (c *SandboxController) UpdateSandbox(w http.ResponseWriter, r *http.Request) error {
	id := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		return domain.InvalidRequestError("invalid sandbox id", nil)
	}

	var req dto.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", nil)
	}

	updates := make(map[string]interface{})
	if req.MemoryLimit != nil {
		updates["MemoryLimit"] = *req.MemoryLimit
	}
	if req.CPULimit != nil {
		updates["CPULimit"] = *req.CPULimit
	}
	if req.PidsLimit != nil {
		updates["PidsLimit"] = *req.PidsLimit
	}
	if req.SessionTimeout != nil {
		updates["SessionTimeout"] = *req.SessionTimeout
	}
	if req.ExecTimeout != nil {
		updates["ExecTimeout"] = *req.ExecTimeout
	}
	if req.NetworkMode != nil {
		updates["NetworkMode"] = *req.NetworkMode
	}

	if len(updates) == 0 {
		return domain.InvalidRequestError("no valid fields to update", nil)
	}

	if err := c.service.UpdateDetails(r.Context(), id, updates); err != nil {
		return err
	}
	return nil
}
