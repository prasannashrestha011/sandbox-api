package dto

import "time"

type CreateWarmPoolRequest struct {
	TemplateID string `json:"templateId" binding:"required"`

	MaxActive int `json:"maxActive" binding:"required"`

	// pool config (WarmPool)
	PoolStatus string `json:"status"` // active/inactive

	// scaling policy
	MinIdleThreshold int `json:"minIdleThreshold" binding:"required"`
	MaxIdleThreshold int `json:"maxIdleThreshold" binding:"required"`

	ScaleUpStep   int `json:"scaleUpStep"`
	ScaleDownStep int `json:"scaleDownStep"`
	CooldownSec   int `json:"cooldownSec"`
}
type WarmPoolResponse struct {
	ID string `json:"id"`

	TemplateID string `json:"templateId"`

	MaxActive int `json:"maxActive"`

	Status string `json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type ScalingPolicyResponse struct {
	ID string `json:"id"`

	WarmPoolID string `json:"warmPoolId"`

	MinIdleThreshold int `json:"minIdleThreshold"`
	MaxIdleThreshold int `json:"maxIdleThreshold"`

	ScaleUpStep   int `json:"scaleUpStep"`
	ScaleDownStep int `json:"scaleDownStep"`

	CooldownSec int `json:"cooldownSec"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateWarmPoolResponse struct {
	Pool   WarmPoolResponse      `json:"pool"`
	Policy ScalingPolicyResponse `json:"policy"`
}
