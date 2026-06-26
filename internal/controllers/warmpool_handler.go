package controllers

import (
	"encoding/json"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/response"
	"main/internal/services"
	"net/http"
)

type WarmPoolHandler struct {
	warmpoolService services.WarmpoolService
}

func NewWarmPoolHandler(svc services.WarmpoolService) *WarmPoolHandler {
	return &WarmPoolHandler{
		warmpoolService: svc,
	}
}

func (h *WarmPoolHandler) CreateWarmPool(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateWarmPoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", nil)
	}
	resp, err := h.warmpoolService.CreateWarmpool(&req)
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "create new warm pool", resp, nil)
	return nil
}
