package mapper

import (
	gorm_model "main/internal/repository/model/lab"
	service_model "main/internal/services/models"
)

// LabToGorm maps a service Lab model to a GORM Lab model
func LabToGorm(l *service_model.Lab) *gorm_model.Lab {
	if l == nil {
		return nil
	}

	chapters := make([]gorm_model.Chapter, len(l.Chapters))
	for i, ch := range l.Chapters {
		chapters[i] = *ChapterToGorm(&ch)
	}

	tags := make([]gorm_model.Tag, len(l.Tags))
	for i, tag := range l.Tags {
		tags[i] = *TagToGorm(&tag)
	}

	return &gorm_model.Lab{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		Lang:        l.Lang,
		ContainerID: l.ContainerID,
		CreatedByID: l.CreatedByID,
		Chapters:    chapters,
		Tags:        tags,
		Difficulty:  l.Difficulty,
		IsPublic:    l.IsPublic,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

// LabFromGorm maps a GORM Lab model to a service Lab model
func LabFromGorm(l *gorm_model.Lab) *service_model.Lab {
	if l == nil {
		return nil
	}

	chapters := make([]service_model.Chapter, len(l.Chapters))
	for i, ch := range l.Chapters {
		chapters[i] = *ChapterFromGorm(&ch)
	}

	tags := make([]service_model.Tag, len(l.Tags))
	for i, tag := range l.Tags {
		tags[i] = *TagFromGorm(&tag)
	}

	return &service_model.Lab{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		Lang:        l.Lang,
		ContainerID: l.ContainerID,
		CreatedByID: l.CreatedByID,
		CreatedBy: service_model.User{
			Username: l.CreatedBy.Username,
			Fullname: l.CreatedBy.Fullname,
		},
		Chapters:   chapters,
		Tags:       tags,
		Difficulty: l.Difficulty,
		IsPublic:   l.IsPublic,
		CreatedAt:  l.CreatedAt,
		UpdatedAt:  l.UpdatedAt,
	}
}
func ChapterToGorm(c *service_model.Chapter) *gorm_model.Chapter {
	if c == nil {
		return nil
	}
	exercises := make([]gorm_model.Exercise, len(c.Exercises))
	for i, ex := range c.Exercises {
		exercises[i] = *ExerciseToGorm(&ex)
	}
	return &gorm_model.Chapter{
		ID:          c.ID,
		LabID:       c.LabID,
		Title:       c.Title,
		Description: c.Description,
		OrderIndex:  c.OrderIndex,
		Exercises:   exercises,
	}
}
func ChapterFromGorm(c *gorm_model.Chapter) *service_model.Chapter {
	if c == nil {
		return nil
	}
	exercises := make([]service_model.Exercise, len(c.Exercises))
	for i, ex := range c.Exercises {
		exercises[i] = *ExerciseFromGorm(&ex)
	}
	return &service_model.Chapter{
		ID:          c.ID,
		LabID:       c.LabID,
		Title:       c.Title,
		Description: c.Description,
		OrderIndex:  c.OrderIndex,
		Exercises:   exercises,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// ExerciseToGorm maps a service Exercise model to a GORM Exercise model
func ExerciseToGorm(e *service_model.Exercise) *gorm_model.Exercise {
	if e == nil {
		return nil
	}
	return &gorm_model.Exercise{
		ID:             e.ID,
		Title:          e.Title,
		ChapterID:      e.ChapterID,
		Description:    e.Description,
		StarterCode:    e.StarterCode,
		ExpectedOutput: e.ExpectedOutput,
		Hints:          e.Hints,
		OrderIndex:     e.OrderIndex,
		Solution:       e.Solution,
		MaxAttempts:    e.MaxAttempts,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

// ExerciseFromGorm maps a GORM Exercise model to a service Exercise model
func ExerciseFromGorm(e *gorm_model.Exercise) *service_model.Exercise {
	if e == nil {
		return nil
	}
	return &service_model.Exercise{
		ID:             e.ID,
		Title:          e.Title,
		Description:    e.Description,
		ChapterID:      e.ChapterID,
		StarterCode:    e.StarterCode,
		ExpectedOutput: e.ExpectedOutput,
		Hints:          e.Hints,
		OrderIndex:     e.OrderIndex,
		Solution:       e.Solution,
		MaxAttempts:    e.MaxAttempts,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

// TagToGorm maps a service Tag model to a GORM Tag model
func TagToGorm(t *service_model.Tag) *gorm_model.Tag {
	if t == nil {
		return nil
	}
	return &gorm_model.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

// TagFromGorm maps a GORM Tag model to a service Tag model
func TagFromGorm(t *gorm_model.Tag) *service_model.Tag {
	if t == nil {
		return nil
	}
	return &service_model.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}
