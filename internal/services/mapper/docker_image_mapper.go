package mapper

import (
	"context"
	request_context "main/internal/context"
	"main/internal/dto"
	"main/internal/services/models"
)

func ToDockerImageModel(ctx context.Context, req *dto.CreateImageRequest) *models.DockerImage {

	if req == nil {
		return nil
	}
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil
	}

	return &models.DockerImage{
		ImageTag:    req.ImageTag,
		Runtime:     req.Runtime,
		CreatedByID: userID.String(),
	}
}

func ToDockerImageResponse(d *models.DockerImage) *dto.DockerImageResponse {
	if d == nil {
		return nil
	}

	return &dto.DockerImageResponse{
		ID:       d.ID,
		ImageTag: d.ImageTag,
		Runtime:  d.Runtime,
	}
}

func ToDockerImageResponses(d []*models.DockerImage) []*dto.DockerImageResponse {
	if d == nil {
		return nil
	}

	responses := make([]*dto.DockerImageResponse, len(d))
	for i, v := range d {
		responses[i] = ToDockerImageResponse(v)
	}
	return responses
}
