package core

import (
	"context"
	"log"

	"main/internal/repository/model"
	"main/internal/sandbox/docker/container"
	"main/internal/sandbox/docker/image"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
)

// SandboxClient defines sandbox lifecycle operations.
type SandboxClient interface {
	Create(ctx context.Context, req *model.Sandbox) error
	Close() error
}

type dockerSandboxClient struct {
	apiClient *client.Client
}

// NewSandboxClient returns a Docker-backed SandboxClient.
func NewSandboxClient() (SandboxClient, error) {
	apiClient, err := client.New(
		client.FromEnv,
		client.WithUserAgent("my-application/1.0.0"),
	)
	if err != nil {
		return nil, err
	}
	return &dockerSandboxClient{apiClient: apiClient}, nil
}

func (c *dockerSandboxClient) Create(ctx context.Context, req *model.Sandbox) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, req.SessionTimeout)
	defer cancel()

	log.Println("Image tag: ", req.Image.ImageTag)
	err := image.PullImage(ctxWithTimeout, c.apiClient, req.Image.ImageTag)
	if err != nil {
		return err
	}
	containerID, err := container.CreateContainer(ctxWithTimeout, c.apiClient, req)
	if err != nil {
		return err
	}
	log.Println("Container ID: ", containerID)
	_, err = c.apiClient.ContainerStart(ctxWithTimeout, containerID, client.ContainerStartOptions{})
	if err != nil {
		return err
	}

	sessionID, _ := uuid.NewUUID()
	req.SessionID = sessionID
	req.ContainerID = containerID
	return nil
}

func (c *dockerSandboxClient) Close() error {
	return c.apiClient.Close()
}
