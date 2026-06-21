package mapper

import (
	"main/internal/enums"
	gorm_model "main/internal/repository/model"
	"main/internal/services/models"
)

func WarmPoolToGorm(wp *models.WarmPool) *gorm_model.WarmPool {
	if wp == nil {
		return nil
	}
	return &gorm_model.WarmPool{
		ID:         wp.ID,
		TemplateID: wp.TemplateID,
		MaxActive:  wp.MaxActive,
		Status:     string(wp.Status),
	}
}

func WarmPoolFromGorm(wp *gorm_model.WarmPool) *models.WarmPool {
	if wp == nil {
		return nil
	}
	return &models.WarmPool{
		ID:         wp.ID,
		TemplateID: wp.TemplateID,
		MaxActive:  wp.MaxActive,
		Status:     enums.PoolStatus(wp.Status),
	}
}

func ScalingPolicyToGorm(sp *models.ScalingPolicy) *gorm_model.ScalingPolicy {
	if sp == nil {
		return nil
	}
	return &gorm_model.ScalingPolicy{
		ID:               sp.ID,
		WarmPoolID:       sp.WarmPoolID,
		MinIdleThreshold: sp.MinIdleThreshold,
		MaxIdleThreshold: sp.MaxIdleThreshold,
		ScaleUpStep:      sp.ScaleUpStep,
		ScaleDownStep:    sp.ScaleDownStep,
		CooldownSec:      sp.CooldownSec,
	}
}

func ScalingPolicyFromGorm(sp *gorm_model.ScalingPolicy) *models.ScalingPolicy {
	if sp == nil {
		return nil
	}
	return &models.ScalingPolicy{
		ID:               sp.ID,
		WarmPoolID:       sp.WarmPoolID,
		MinIdleThreshold: sp.MinIdleThreshold,
		MaxIdleThreshold: sp.MaxIdleThreshold,
		ScaleUpStep:      sp.ScaleUpStep,
		ScaleDownStep:    sp.ScaleDownStep,
		CooldownSec:      sp.CooldownSec,
	}
}
