package models

import (
	"main/internal/enums"
	"time"
)

type WarmPool struct {
	ID string

	TemplateID string

	MaxActive int

	Status enums.PoolStatus
}

type ScalingPolicy struct {
	ID string

	WarmPoolID string

	MinIdleThreshold int
	MaxIdleThreshold int

	ScaleUpStep   int
	ScaleDownStep int

	CooldownSec int
}
type PoolRuntimeState struct {
	WarmPoolID string

	CurrentIdle   int
	CurrentActive int

	LastScaledAt time.Time
}
