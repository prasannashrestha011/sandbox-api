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

func SessionToGorm(s *models.SandboxSession) *gormodel.SandboxInstance {
	return &gormodel.SandboxInstance{
		ID:          s.ID,
		UserID:      s.UserID,
		TemplateID:  s.TemplateID,
		ContainerID: s.ContainerID,
		Status:      string(s.Status),
		Lang:        s.Lang,
	}
}

func SessionFromGorm(m *gormodel.SandboxInstance) *models.SandboxSession {
	return &models.SandboxSession{
		ID:          m.ID,
		UserID:      m.UserID,
		TemplateID:  m.TemplateID,
		ContainerID: m.ContainerID,
		Status:      enums.SandboxState(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Lang:        m.Lang,
	}
}
