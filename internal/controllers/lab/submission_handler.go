package controllers

import (
	"encoding/json"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/pkg"
	"main/internal/response"
	lab_services "main/internal/services/lab"
	"net/http"

	"github.com/google/uuid"
)

type SubmissionHandler struct {
	submissionService lab_services.SubmissionService
}

func NewSubmissionHandler(submissionService lab_services.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
	}
}

func (h *SubmissionHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) error {
	exerciseID := pkg.ExtractParam(r, "exerciseId")
	if _, err := uuid.Parse(exerciseID); err != nil {
		return domain.InvalidRequestError("invalid exercise id", nil)
	}
	var req dto.SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid submission request", nil)
	}
	resp, err := h.submissionService.CreateSubmission(r.Context(), exerciseID, &req)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusCreated, "solution submitted", resp, nil)
	return nil
}
func (h *SubmissionHandler) GetSubmission(w http.ResponseWriter, r *http.Request) error {
	id := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		return domain.InvalidRequestError("invalid submission id", nil)
	}
	resp, err := h.submissionService.GetSubmission(r.Context(), id)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "submission details", resp, nil)
	return nil
}

func (h *SubmissionHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) error {
	exerciseID := pkg.ExtractParam(r, "exerciseId")
	if _, err := uuid.Parse(exerciseID); err != nil {
		return domain.InvalidRequestError("invalid exercise id", nil)
	}
	resp, err := h.submissionService.ListSubmissions(r.Context(), exerciseID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "list of submissions for the exercise", resp, nil)
	return nil
}

func (h *SubmissionHandler) UpdateSubmission(w http.ResponseWriter, r *http.Request) error {
	id := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		return domain.InvalidRequestError("invalid submission id", nil)
	}
	var req dto.SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid submission request", nil)
	}
	resp, err := h.submissionService.UpdateSubmission(r.Context(), id, &req)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "submission updated", resp, nil)
	return nil
}

func (h *SubmissionHandler) DeleteSubmission(w http.ResponseWriter, r *http.Request) error {
	id := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		return domain.InvalidRequestError("invalid submission id", nil)
	}
	err := h.submissionService.DeleteSubmission(r.Context(), id)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "submission deleted", nil, nil)
	return nil
}
