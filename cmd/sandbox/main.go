package main

import (
	"context"

	app "main/internal"
	"main/internal/database"
	"main/internal/sandbox/core"
)

func main() {
	// docker client
	// client.FromEnv == reads docker connection string from environment

	ctx := context.Background()

	apiClient, err := core.NewSandboxClient()
	if err != nil {
		panic(err)
	}

	defer apiClient.Close()

	db, err := database.ConnectFromEnv(ctx)
	if err != nil {
		panic(err)
	}

	application, err := app.New(db, apiClient)
	if err != nil {
		panic(err)
	}
	_ = application

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT)
	//
	// go func() {
	// 	<-sigChan
	// 	defer apiClient.ImageRemove(
	// 		ctx,
	// 		imageID,
	// 		client.ImageRemoveOptions{
	// 			Force:         true,
	// 			PruneChildren: true,
	// 		},
	// 	)
	// 	defer apiClient.ContainerRemove(
	// 		ctx,
	// 		resp.ContainerID,
	// 		client.ContainerRemoveOptions{
	// 			Force: true,
	// 		},
	// 	)
	// }()
	//
	// defer apiClient.ImageRemove(
	// 	ctx,
	// 	imageID,
	// 	client.ImageRemoveOptions{
	// 		Force:         true,
	// 		PruneChildren: true,
	// 	},
	// )
	//
	// defer apiClient.ContainerRemove(ctx, resp.ContainerID, client.ContainerRemoveOptions{
	// 	Force: true,
	// })

}
