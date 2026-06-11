package container

import (
	"context"
	"fmt"
	"main/internal/pkg"
	"main/internal/services/models"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func CreateContainer(ctx context.Context, apiClient *client.Client, req *models.SandboxTemplate) (containerID string, containerName string, err error) {
	containerName, _ = pkg.GenRandomString(12)
	resp, err := apiClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image:     req.Image.ImageTag,
			Cmd:       []string{"sleep", "infinity"},
			Tty:       true,
			OpenStdin: true,
			StdinOnce: false,
			Labels: map[string]string{
				"app": "sandbox",
			},
		},
		HostConfig: &container.HostConfig{
			NetworkMode: container.NetworkMode(req.NetworkMode),
			Resources: container.Resources{
				Memory:    req.MemoryLimit,
				NanoCPUs:  req.CPULimit,
				PidsLimit: &[]int64{req.PidsLimit}[0],
			},
		},
		NetworkingConfig: nil,
		Platform:         nil,
		Name:             fmt.Sprintf("sandbox_%s", containerName),
	})
	if err != nil {
		return "", "", err
	}
	return resp.ID, containerName, nil
}
