package controllers

import (
	"encoding/json"
	"errors"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/response"
	"main/internal/services"
	"net/http"
)

type DockerImageController struct {
	service services.DockerImageService
}

func NewDockerImageController(service services.DockerImageService) *DockerImageController {
	return &DockerImageController{service: service}
}

func (c *DockerImageController) CreateImage(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domain.InvalidRequestError("invalid request body", err)
	}
	if err := req.Validate(); err != nil {
		var v *dto.ValidationErrors
		if errors.As(err, &v) {
			return domain.ValidationError(err, v.Violations)
		}
		return domain.ValidationError(err, nil)

	}

	createdImage, err := c.service.CreateImage(r.Context(), &req)
	if err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusCreated, "docker image record created successfully", createdImage, nil)
	return nil
}
func (c *DockerImageController) ListImages(w http.ResponseWriter, r *http.Request) error {
	images, err := c.service.ListImages(r.Context())
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "docker images found", images, nil)
	return nil
}
