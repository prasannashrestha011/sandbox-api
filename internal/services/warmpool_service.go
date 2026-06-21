package services

import (
	"main/internal/dto"
	postgres_error "main/internal/infra/postgres"
	repository "main/internal/repository"
	"main/internal/services/mapper"
)

type WarmpoolService interface {
	CreateWarmpool(req *dto.CreateWarmPoolRequest) (*dto.WarmPoolResponse, error)
}
type warmpoolService struct {
	warmpoolRepo repository.WarmPoolRepository
}

func NewWarmpoolService(warmpoolRepo repository.WarmPoolRepository) WarmpoolService {
	return &warmpoolService{
		warmpoolRepo: warmpoolRepo,
	}
}
func (w *warmpoolService) CreateWarmpool(req *dto.CreateWarmPoolRequest) (*dto.WarmPoolResponse, error) {
	warmPool := mapper.ToWarmPoolModel(req)
	createdWarmPool, err := w.warmpoolRepo.CreateWarmPool(warmPool)
	if err != nil {
		return nil, postgres_error.MapError(err, "create warmpool", "warmpool")
	}
	return mapper.ToWarmPoolResponse(createdWarmPool), nil
}
