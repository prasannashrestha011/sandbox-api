package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"main/internal/domain"
	"main/internal/dto"
	"main/internal/pkg"
	"main/internal/response"
	lab_services "main/internal/services/lab"
)

type LabController struct {
	service lab_services.LabService
}

func NewLabController(service lab_services.LabService) *LabController {
	return &LabController{service: service}
}

func (c *LabController) CreateLab(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateLabRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request JSON body", err)
	}

	if err := req.Validate(); err != nil {
		var v *dto.ValidationErrors
		if errors.As(err, &v) {
			return domain.ValidationError(err, v.Violations)
		}
		return domain.ValidationError(err, nil)
	}
	resp, err := c.service.CreateLab(r.Context(), &req)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusCreated, "lab created successfully", resp, nil)
	return nil
}

func (c *LabController) GetLabByID(w http.ResponseWriter, r *http.Request) error {
	idStr := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(idStr); err != nil {
		return domain.InvalidRequestError("invalid lab id", nil)
	}

	resp, err := c.service.GetLabByID(r.Context(), idStr)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusOK, "lab retrieved successfully", resp, nil)
	return nil
}

func (c *LabController) DeleteLab(w http.ResponseWriter, r *http.Request) error {
	idStr := pkg.ExtractParam(r, "id")
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
