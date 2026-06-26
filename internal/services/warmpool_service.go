package services

import (
	"context"
	"log"
	"main/internal/dto"
	postgres_error "main/internal/infra/postgres"
	"main/internal/jobs/warmpool"
	repository "main/internal/repository"
	"main/internal/services/mapper"

	"github.com/hibiken/asynq"
)

type WarmpoolService interface {
	CreateWarmpool(req *dto.CreateWarmPoolRequest) (*dto.WarmPoolResponse, error)
	Execute(ctx context.Context, req *dto.SandboxExecReq) (*dto.SandboxExecResponse, error)
}
type warmpoolService struct {
	warmpoolRepo   repository.WarmPoolRepository
	sandboxService SandboxInstanceService
	asyncClient    *asynq.Client
}

func NewWarmpoolService(warmpoolRepo repository.WarmPoolRepository, sandboxService SandboxInstanceService, asyncClient *asynq.Client) WarmpoolService {
	return &warmpoolService{
		warmpoolRepo:   warmpoolRepo,
		sandboxService: sandboxService,
		asyncClient:    asyncClient,
	}
}
func (w *warmpoolService) CreateWarmpool(req *dto.CreateWarmPoolRequest) (*dto.WarmPoolResponse, error) {
	warmPool := mapper.ToWarmPoolModel(req)
	createdWarmPool, err := w.warmpoolRepo.CreateWarmPool(warmPool)
	if err != nil {
		return nil, postgres_error.MapError(err, "create warmpool", "warmpool")
	}

	scalingPolicy := mapper.ToScalingPolicyModel(req, warmPool.ID)
	payload := &dto.SandboxProvisionPayload{
		TemplateID: createdWarmPool.TemplateID,
		PoolID:     createdWarmPool.ID,
	}
	task, _ := warmpool.NewSandboxProvisionTask(payload)
	for i := 0; i < scalingPolicy.MinIdleThreshold; i++ {
		log.Println("Enqueuing task to provision sandbox for warm pool:", createdWarmPool.ID, "Template:", createdWarmPool.TemplateID)
		_, err = w.asyncClient.Enqueue(task)
		if err != nil {
			return nil, err
		}
	}

	return mapper.ToWarmPoolResponse(createdWarmPool), nil
}

func (w *warmpoolService) Execute(ctx context.Context, req *dto.SandboxExecReq) (*dto.SandboxExecResponse, error) {
	sandbox, err := w.sandboxService.AcquireInstance(ctx, req.Lang)
	if err != nil {
		return nil, err
	}
	resp, err := w.sandboxService.ExecuteCommand(ctx, sandbox.ContainerID, req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
