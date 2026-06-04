package controllers

import (
	"encoding/json"
	"log"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/pkg"
	"main/internal/response"
	lab_services "main/internal/services/lab"
	"net/http"

	"github.com/google/uuid"
)

type ChapterController struct {
	chapterService lab_services.ChapterService
}

func NewChapterController(chapterService lab_services.ChapterService) *ChapterController {
	return &ChapterController{chapterService: chapterService}
}

func (h *ChapterController) CreateChapter(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateChapterRequest
	labID := pkg.ExtractParam(r, "id")
	if _, err := uuid.Parse(labID); err != nil {
		return domain.InvalidRequestError("invalid lab id", nil)
	}

	log.Println("labID: ", labID)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request JSON body", err)
	}

	resp, err := h.chapterService.CreateChapter(r.Context(), &req, labID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusCreated, "chapter created successfully", resp, nil)
	return nil
}

func (h *ChapterController) GetChaptersByLabID(w http.ResponseWriter, r *http.Request) error {
	labID := pkg.ExtractParam(r, "labId")
	if _, err := uuid.Parse(labID); err != nil {
		return domain.InvalidRequestError("invalid lab id", nil)
	}

	resp, err := h.chapterService.GetChaptersByLabID(r.Context(), labID)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "chapters retrieved successfully", resp, nil)
	return nil
}

func (h *ChapterController) UpdateChapter(w http.ResponseWriter, r *http.Request) error {
	var req dto.UpdateChapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request JSON body", err)
	}

	resp, err := h.chapterService.UpdateChapter(r.Context(), &req)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "chapter updated successfully", resp, nil)
	return nil
}

func (h *ChapterController) DeleteChapter(w http.ResponseWriter, r *http.Request) error {
	chapterID := pkg.ExtractParam(r, "chapterId")
	if _, err := uuid.Parse(chapterID); err != nil {
		return domain.InvalidRequestError("invalid chapter id", nil)
	}

	if err := h.chapterService.DeleteChapter(r.Context(), chapterID); err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "chapter deleted successfully", nil, nil)
	return nil
}
