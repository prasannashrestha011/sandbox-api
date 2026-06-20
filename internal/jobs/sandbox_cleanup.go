package jobs

import (
	"context"
	"encoding/json"
	"log"
	"main/internal/database"
	"main/internal/dto"
	"main/internal/repository/model"
	"time"

	"github.com/hibiken/asynq"
	"github.com/moby/moby/client"
	"gorm.io/gorm"
)

const TypeSandboxCleanup = "sandbox:cleanup"

func NewSandboxCleanupTask(payload *dto.SandboxCleanupPayload) (*asynq.Task, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSandboxCleanup, jsonPayload), nil
}

var dockerClient *client.Client
var db *gorm.DB

func init() {
	var err error
	dockerClient, err = client.New(
		client.FromEnv,
		client.WithUserAgent("my-application/1.0.0"),
	)
	if err != nil {
		log.Println("Error initialzing docker client ", err.Error())
	}
	ctx := context.Background()
	db, err = database.ConnectFromEnv(ctx)
	if err != nil {
		panic(err)
	}

}

func HandleSandboxCleanup(ctx context.Context, t *asynq.Task) error {
	var p dto.SandboxCleanupPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	log.Println("Starting cleanup : ", p.ContainerID)
	log.Println("SessionID: ", p.SessionID)

	filterArgs := client.Filters{"label": {"app=sandbox": true}, "id": {p.ContainerID: true}}
	result, err := dockerClient.ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		log.Println("Error listing containers: ", err)
		return err
	}
	if len(result.Items) == 0 {
		log.Println("No container found with ID: ", p.ContainerID)
		return nil
	}
	ctr := result.Items[0]
	imageID := ctr.ImageID

	removedResult, err := dockerClient.ContainerRemove(ctx, ctr.ID, client.ContainerRemoveOptions{Force: true})
	if err != nil {
		log.Println("Error removing container: ", err)
		return err
	}
	log.Println("Container removed result", removedResult)
	_, err = dockerClient.ImageRemove(ctx, imageID, client.ImageRemoveOptions{Force: true})
	if err != nil {
		log.Println("Error removing image: ", err)
		return err
	}
	// after successful container removal
	err = db.Model(&model.SandboxInstance{}).
		Where("container_id = ?", p.ContainerID).
		Update("status", "inactive").Error

	return nil
}
