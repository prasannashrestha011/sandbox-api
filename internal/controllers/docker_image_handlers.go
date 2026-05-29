package controllers

import (
	"encoding/json"
	request_context "main/internal/context"
	"main/internal/services"
	"main/internal/types"
	"net/http"
)

type DockerImageController struct {
	service services.DockerImageService
}

func NewDockerImageController(service services.DockerImageService) *DockerImageController {
	return &DockerImageController{service: service}
}

func (c *DockerImageController) CreateImage(w http.ResponseWriter, r *http.Request) {
	var req types.CreateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	userID, ok := request_context.UserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	if err := c.service.CreateImage(req.ImageTag, userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create docker image record"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "docker image record created successfully"})
}
