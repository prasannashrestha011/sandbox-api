package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	sandbox_exec "main/actions/exec"
	sandbox_executor "main/actions/executor"
	sandbox_image "main/actions/image"
	sandbox_client "main/client"
	sandbox_request "main/types"

	"github.com/moby/moby/client"
)

func main() {
	// docker client
	// client.FromEnv == reads docker connection string from environment

	ctx := context.Background()

	apiClient, err := sandbox_client.NewSandboxClient()
	if err != nil {
		panic(err)
	}

	defer apiClient.Close()

	// pull the image
	base := sandbox_executor.LangCommands["javascript"]
	imageID := sandbox_image.LoadImage("javascript")

	// creating a container
	req := &sandbox_request.CreateRequest{
		UserID:         "8080080",
		Environment:    "javascript",
		ImageID:        imageID,
		MemoryLimit:    64 * 1024 * 1024,
		CPULimit:       500000000,
		PidsLimit:      12,
		NetWorkMode:    "none",
		CreatedAt:      time.Now(),
		SessionTimeout: time.Second * 600,
		ExecTimeout:    time.Second * 60,
	}

	resp, err := sandbox_client.CreateNewSandBox(apiClient, req)
	if err != nil {
		panic(resp)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		defer apiClient.ImageRemove(
			ctx,
			imageID,
			client.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			},
		)
		defer apiClient.ContainerRemove(
			ctx,
			resp.ContainerID,
			client.ContainerRemoveOptions{
				Force: true,
			},
		)
	}()

	defer apiClient.ImageRemove(
		ctx,
		imageID,
		client.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		},
	)

	defer apiClient.ContainerRemove(ctx, resp.ContainerID, client.ContainerRemoveOptions{
		Force: true,
	})

	// start the container

	// test commands
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		cmd = strings.ToLower(cmd)
		if cmd == "exit" {
			break
		}
		dockerCmd := append(base, cmd)
		sandbox_exec.ExecCreate(ctx, apiClient, resp.ContainerID, dockerCmd)
	}
}
