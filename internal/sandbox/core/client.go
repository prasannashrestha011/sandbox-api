package core

import (
	"context"
	"log"

	"main/internal/repository/model"
	"main/internal/sandbox/docker/container"
	sb_executil "main/internal/sandbox/docker/executil"
	"main/internal/sandbox/docker/image"

	"github.com/google/uuid"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
)

// SandboxClient defines sandbox lifecycle operations.
type SandboxClient interface {
	Create(ctx context.Context, req *model.Sandbox) error
	ExecuteCode(ctx context.Context, containerID string, cmd []string) (string, error)
	CleanUp(ctx context.Context, handler func(containerID []string)) error
	ListenContainerEvents(ctx context.Context, handler func(containerID string)) error
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
func (c *dockerSandboxClient) ExecuteCode(ctx context.Context, containerID string, cmd []string) (string, error) {
	return sb_executil.ExecCreate(ctx, c.apiClient, containerID, cmd)
}

func (c *dockerSandboxClient) CleanUp(ctx context.Context, handler func(containerID []string)) error {
	result, err := c.apiClient.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		log.Println("Clean up func: ", err.Error())
		return err
	}

	imageIDs := make(map[string]struct{})
	containerIDs := make([]string, 0)
	//Remove containers labeled with "app=sandbox" and collect their image IDs for later cleanup
	for _, ctr := range result.Items {

		if ctr.Labels["app"] == "sandbox" {
			imageIDs[ctr.ImageID] = struct{}{}
			containerIDs = append(containerIDs, ctr.ID)
			log.Printf("CleanUp: found container %s with image %s", ctr.ID, ctr.ImageID)
			_, err := c.apiClient.ContainerRemove(ctx, ctr.ID, client.ContainerRemoveOptions{
				Force: true,
			})
			if err != nil {
				log.Println("Clean up func: ", err.Error())
				return err
			}
		}

	}
	log.Printf("CleanUp: found %d total containers", len(result.Items))
	// Remove images that are no longer in use
	for imageID := range imageIDs {
		log.Printf("CleanUp: removing image %s", imageID)
		_, err := c.apiClient.ImageRemove(ctx, imageID, client.ImageRemoveOptions{
			Force: true,
		})
		if err != nil {
			log.Printf("Clean up func: %v", err.Error())
			return err
		}
	}
	handler(containerIDs)
	return nil
}

// listens for container events and calls the handler when a container dies or stops. This is used to clean up resources associated with the container.
func (c *dockerSandboxClient) ListenContainerEvents(ctx context.Context, handler func(containerID string)) error {
	eventResult := c.apiClient.Events(ctx, client.EventsListOptions{})
	msgs, errs := eventResult.Messages, eventResult.Err
	for {
		select {
		case msg := <-msgs:
			if msg.Type == events.ContainerEventType && (msg.Action == "die" || msg.Action == "stop") {
				handler(msg.Actor.ID)
			}
		case err := <-errs:
			return err

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *dockerSandboxClient) Close() error {
	return c.apiClient.Close()
}
