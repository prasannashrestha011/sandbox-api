package warmpool

import (
	"context"
	"encoding/json"
	"log"
	"main/internal/database"
	"main/internal/dto"
	"main/internal/enums"
	"main/internal/repository"
	"main/internal/sandbox/core"
	"main/internal/services/models"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const TypeWarmPool = "warmpool:provision"

func NewWarmPoolCreateTask(payload *dto.CreateWarmPoolRequest) (*asynq.Task, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWarmPool, jsonPayload), nil
}

var db *gorm.DB
var templateRepo repository.SandboxTemplateRepository
var sandboxRepo repository.SandboxInstanceRepository

var dockerClient core.SandboxClient

func init() {
	var err error
	dockerClient, err = core.NewSandboxClient()
	ctx := context.Background()
	db, err = database.ConnectFromEnv(ctx)
	if err != nil {
		panic(err)
	}
	templateRepo = repository.NewSandboxTemplateRepository(db)
	sandboxRepo = repository.NewSandboxInstanceRepository(db)

}
func NewSandboxProvisionTask(payload *dto.SandboxProvisionPayload) (*asynq.Task, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWarmPool, jsonPayload), nil
}

func HandleSandboxProvision(ctx context.Context, t *asynq.Task) error {
	var p dto.SandboxProvisionPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	// 2. Find the language/sandbox template rules
	template, err := templateRepo.FindByID(ctx, p.TemplateID)
	if err != nil {
		log.Println("Error finding sandbox template: ", err.Error())
		return err
	}

	// 3. Create the container using the Moby SDK
	// Assumes dockerClient.Create returns the Container ID and networking port/metadata
	containerID, _, err := dockerClient.Create(ctx, template)
	if err != nil {
		log.Println("Error creating sandbox container: ", err.Error())
		return err // Asynq will automatically retry if you return the error
	}

	// 4. Track this single instance in your database
	sandbox := &models.SandboxInstance{
		PoolID:      p.PoolID,
		Status:      enums.StateActive, // Mark it as ready for action
		ContainerID: containerID,
		Lang:        template.Lang,
		LastUsedAt:  time.Now(),
	}
	_, err = sandboxRepo.Create(ctx, sandbox)
	if err != nil {
		log.Println("Error saving sandbox instance to DB: ", err.Error())
		return err
	}
	log.Println("Sandbox container created and tracked successfully: ", containerID)
	return nil
}
