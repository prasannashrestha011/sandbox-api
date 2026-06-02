package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"main/internal/domain"
	"main/internal/dto"
	"main/internal/services"
)

type LabController struct {
	service services.LabService
}

func NewLabController(service services.LabService) *LabController {
	return &LabController{service: service}
}

func (c *LabController) CreateLab(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateLabRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request JSON body", err)
	}

	resp, err := c.service.CreateLab(r.Context(), &req)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, resp)
	return nil
}

func (c *LabController) GetLabByID(w http.ResponseWriter, r *http.Request) error {
	idStr := extractParam(r, "id")
	if _, err := uuid.Parse(idStr); err != nil {
		return domain.InvalidRequestError("invalid lab id", nil)
	}

	resp, err := c.service.GetLabByID(r.Context(), idStr)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, resp)
	return nil
}

func (c *LabController) DeleteLab(w http.ResponseWriter, r *http.Request) error {
	idStr := extractParam(r, "id")
	if _, err := uuid.Parse(idStr); err != nil {
		return domain.InvalidRequestError("invalid lab id", nil)
	}

	err := c.service.DeleteLab(r.Context(), idStr)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
