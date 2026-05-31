package controllers

import (
	"encoding/json"
	"errors"
	request_context "main/internal/context"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/response"
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

func (c *DockerImageController) CreateImage(w http.ResponseWriter, r *http.Request) error {
	var req types.CreateImageRequest
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

	userID, ok := request_context.UserID(r.Context())
	if !ok {
		return domain.InvalidRequestError("missing user id", nil)
	}

	if err := c.service.CreateImage(req.ImageTag, userID.String()); err != nil {
		return err
	}

	response.WriteJSON(w, r, http.StatusCreated, "docker image record created successfully", nil, nil)
	return nil
}
func (c *DockerImageController) ListImages(w http.ResponseWriter, r *http.Request) error {
	images, err := c.service.ListImages()
	if err != nil {
		return err
	}
	response.WriteJSON(w, r, http.StatusOK, "docker images found", images, nil)
	return nil
}
