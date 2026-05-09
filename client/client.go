package sandbox_client

import (
	"context"
	"log"
	"time"

	sandbox_container "main/actions/container"
	sandbox_image "main/actions/image"
	sandbox_request "main/types"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
)

func NewSandboxClient() (*client.Client, error) {
	apiClient, err := client.New(
		client.FromEnv,
		client.WithUserAgent("my-application/1.0.0"),
	)
	if err != nil {
		return nil, err
	}
	return apiClient, nil
}

func CreateNewSandBox(apiClient *client.Client, req *sandbox_request.CreateRequest) (*sandbox_request.CreateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), req.SessionTimeout)
	defer cancel()
	sandbox_image.PullImage(ctx, apiClient, req.ImageID)
	containerID, err := sandbox_container.CreateContainer(ctx, apiClient, req)
	if err != nil {
		return nil, err
	}
	log.Println("Container ID: ", containerID)

	log.Println("Starting the container")
	_, err = apiClient.ContainerStart(ctx, containerID, client.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	sessionID, _ := uuid.NewUUID()
	return &sandbox_request.CreateResponse{
		ContainerID: containerID,
		SessionID:   sessionID,
		Status:      "active",
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(req.SessionTimeout),
	}, nil
}
