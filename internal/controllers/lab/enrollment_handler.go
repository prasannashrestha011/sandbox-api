package controllers

import (
	"encoding/json"
	request_context "main/internal/context"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/pkg"
	"main/internal/response"
	lab_services "main/internal/services/lab"
	"net/http"

	"github.com/google/uuid"
)

type EnrollmentController struct {
	enrollmentService lab_services.EnrollmentService
}

func NewEnrollmentController(enrollmentService lab_services.EnrollmentService) *EnrollmentController {
	return &EnrollmentController{enrollmentService: enrollmentService}
}

func (h *EnrollmentController) EnrollUserToLab(w http.ResponseWriter, r *http.Request) error {
	var req dto.EnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request JSON body", err)
	}
	err := h.enrollmentService.EnrollUserToLab(r.Context(), &req)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusCreated, "user enrolled to lab successfully", nil, nil)
	return nil
}

func (h *EnrollmentController) GetEnrollment(w http.ResponseWriter, r *http.Request) error {
	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.InvalidRequestError("user id not found in context", nil)
	}
	labID := pkg.ExtractParam(r, "labId")
	if err := validateUUID(userID.String(), "user"); err != nil {
		return err
	}
	if err := validateUUID(labID, "lab"); err != nil {
		return err
	}
	resp, err := h.enrollmentService.GetEnrollment(r.Context(), userID.String(), labID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "user enrollment", resp, nil)
	return nil
}

func (h *EnrollmentController) GetUserEnrollments(w http.ResponseWriter, r *http.Request) error {
	userID := pkg.ExtractParam(r, "userId")
	if err := validateUUID(userID, "user"); err != nil {
		return err
	}
	resp, err := h.enrollmentService.GetUserEnrollments(r.Context(), userID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "user enrollments", resp, nil)
	return nil
}

func (h *EnrollmentController) DeleteEnrollment(w http.ResponseWriter, r *http.Request) error {
	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.InvalidRequestError("user id not found in context", nil)
	}
	labID := pkg.ExtractParam(r, "labId")
	if err := validateUUID(userID.String(), "user"); err != nil {
		return err
	}
	if err := validateUUID(labID, "lab"); err != nil {
		return err
	}
	err := h.enrollmentService.DeleteEnrollment(r.Context(), userID.String(), labID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "enrollment deleted successfully", nil, nil)
	return nil
}

func validateUUID(id string, idType string) error {
	if _, err := uuid.Parse(id); err != nil {
		return domain.InvalidRequestError("invalid "+idType+" id", nil)
	}
	return nil
}
