package image

import (
	"context"
	"io"
	"log"

	"github.com/moby/moby/client"
)

func PullImage(ctx context.Context, apiClient *client.Client, imageTag string) error {
	_, err := apiClient.ImageInspect(ctx, imageTag)
	if err != nil {
		log.Println("Pulling image: ", imageTag)
		reader, err := apiClient.ImagePull(ctx, imageTag, client.ImagePullOptions{})
		if err != nil {
			log.Println("docker: image tag not found")
			return err
		}
		io.Copy(io.Discard, reader)

	}
	return nil
}
