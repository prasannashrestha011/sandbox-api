package container

import (
	"context"
	"fmt"
	"main/internal/pkg"
	"main/internal/repository/model"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func CreateContainer(ctx context.Context, apiClient *client.Client, req *model.Sandbox) (string, error) {
	req.ContainerName, _ = pkg.GenRandomString(12)
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
		Name:             fmt.Sprintf("sandbox_%s", req.ContainerName),
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
