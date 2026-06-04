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

type ExerciseHandler struct {
	exerciseService lab_services.ExerciseService
}

func NewExerciseHandler(exerciseService lab_services.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
	}
}

func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) error {
	chapterID := pkg.ExtractParam(r, "chapterId")
	if _, err := uuid.Parse(chapterID); err != nil {
		return domain.InvalidRequestError("invalid chapterID", nil)
	}
	var req dto.CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", err)
	}
	resp, err := h.exerciseService.CreateExercise(r.Context(), &req, chapterID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusCreated, "exercise created successfully", resp, nil)
	return nil
}

func (h *ExerciseHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request) error {
	exerciseID := pkg.ExtractParam(r, "exerciseId")
	if _, err := uuid.Parse(exerciseID); err != nil {
		return domain.InvalidRequestError("invalid exerciseID", nil)
	}
	resp, err := h.exerciseService.GetExerciseByID(r.Context(), exerciseID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "exercise fetched successfully", resp, nil)
	return nil
}

func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) error {
	exerciseID := pkg.ExtractParam(r, "exerciseId")
	if _, err := uuid.Parse(exerciseID); err != nil {
		return domain.InvalidRequestError("invalid exerciseID", nil)
	}
	var req dto.UpdateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", err)
	}
	err := h.exerciseService.UpdateExercise(r.Context(), exerciseID, &req)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "exercise updated successfully", nil, nil)
	return nil
}

func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) error {
	exerciseID := pkg.ExtractParam(r, "exerciseId")
	if _, err := uuid.Parse(exerciseID); err != nil {
		return domain.InvalidRequestError("invalid exerciseID", nil)
	}
	err := h.exerciseService.DeleteExercise(r.Context(), exerciseID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "exercise deleted successfully", nil, nil)
	return nil
}

func (h *ExerciseHandler) ListExercisesByChapterID(w http.ResponseWriter, r *http.Request) error {
	chapterID := pkg.ExtractParam(r, "chapterId")
	if _, err := uuid.Parse(chapterID); err != nil {
		return domain.InvalidRequestError("invalid chapterID", nil)
	}
	resp, err := h.exerciseService.ListExercisesByChapterID(r.Context(), chapterID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "exercises fetched successfully", resp, nil)
	return nil
}
