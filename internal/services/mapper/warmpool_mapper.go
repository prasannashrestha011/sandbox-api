package mapper

import (
	"main/internal/dto"
	"main/internal/enums"
	"main/internal/services/models"
)

func ToWarmPoolModel(dto *dto.CreateWarmPoolRequest) *models.WarmPool {
	return &models.WarmPool{
		TemplateID: dto.TemplateID,
		MaxActive:  dto.MaxActive,
		Status:     enums.PoolStatus(dto.PoolStatus),
	}
}
func ToScalingPolicyModel(dto *dto.CreateWarmPoolRequest, warmPoolID string) *models.ScalingPolicy {
	return &models.ScalingPolicy{
		WarmPoolID:       warmPoolID,
		MinIdleThreshold: dto.MinIdleThreshold,
		MaxIdleThreshold: dto.MaxIdleThreshold,
		ScaleUpStep:      dto.ScaleUpStep,
		ScaleDownStep:    dto.ScaleDownStep,
		CooldownSec:      dto.CooldownSec,
	}
}

func ToWarmPoolResponse(model *models.WarmPool) *dto.WarmPoolResponse {
	return &dto.WarmPoolResponse{
		ID:         model.ID,
		TemplateID: model.TemplateID,
		MaxActive:  model.MaxActive,
		Status:     string(model.Status),
	}
}
