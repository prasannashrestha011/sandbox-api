package mapper

import (
	"main/internal/enums"
	gormodel "main/internal/repository/model"
	"main/internal/services/models"
)

func TemplateToGorm(t *models.SandboxTemplate) *gormodel.SandboxTemplate {
	return &gormodel.SandboxTemplate{
		ID:     t.ID,
		UserID: t.UserID,
		Lang:   t.Lang,

		ImageID:     t.Image.ID,
		MemoryLimit: t.MemoryLimit,
		CPULimit:    t.CPULimit,
		PidsLimit:   t.PidsLimit,

		SessionTimeout: t.SessionTimeout,
		ExecTimeout:    t.ExecTimeout,
		NetworkMode:    t.NetworkMode,
	}
}
func TemplateFromGorm(m *gormodel.SandboxTemplate) *models.SandboxTemplate {
	return &models.SandboxTemplate{
		ID:     m.ID,
		UserID: m.UserID,
		Lang:   m.Lang,

		Image: models.DockerImage{
			ID:       m.Image.ID,
			ImageTag: m.Image.ImageTag,
		},

		MemoryLimit: m.MemoryLimit,
		CPULimit:    m.CPULimit,
		PidsLimit:   m.PidsLimit,

		SessionTimeout: m.SessionTimeout,
		ExecTimeout:    m.ExecTimeout,
		NetworkMode:    m.NetworkMode,
	}
}

func TemplateListFromGorm(list []gormodel.SandboxTemplate) []models.SandboxTemplate {
	result := make([]models.SandboxTemplate, len(list))
	for i, m := range list {
		result[i] = *TemplateFromGorm(&m)
	}
	return result
}

func InstanceToGorm(s *models.SandboxInstance) *gormodel.SandboxInstance {
	return &gormodel.SandboxInstance{
		ID:          s.ID,
		ContainerID: s.ContainerID,
		Status:      string(s.Status),
		Lang:        s.Lang,
		PoolID:      s.PoolID,
	}
}

func InstanceFromGorm(m *gormodel.SandboxInstance) *models.SandboxInstance {
	return &models.SandboxInstance{
		ID:          m.ID,
		ContainerID: m.ContainerID,
		Status:      enums.SandboxState(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Lang:        m.Lang,
	}
}
