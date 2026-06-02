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

	exercises := make([]gorm_model.Exercise, len(l.Exercises))
	for i, ex := range l.Exercises {
		exercises[i] = *ExerciseToGorm(&ex)
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
		Exercises:   exercises,
		ContainerID: l.ContainerID,
		CreatedByID: l.CreatedByID,
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

	exercises := make([]service_model.Exercise, len(l.Exercises))
	for i, ex := range l.Exercises {
		exercises[i] = *ExerciseFromGorm(&ex)
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
		Exercises:   exercises,
		ContainerID: l.ContainerID,
		CreatedByID: l.CreatedByID,
		CreatedBy: service_model.User{
			Username: l.CreatedBy.Username,
			Fullname: l.CreatedBy.Fullname,
		},
		Tags:       tags,
		Difficulty: l.Difficulty,
		IsPublic:   l.IsPublic,
		CreatedAt:  l.CreatedAt,
		UpdatedAt:  l.UpdatedAt,
	}
}

// ExerciseToGorm maps a service Exercise model to a GORM Exercise model
func ExerciseToGorm(e *service_model.Exercise) *gorm_model.Exercise {
	if e == nil {
		return nil
	}
	return &gorm_model.Exercise{
		ID:             e.ID,
		LabID:          e.LabID,
		Title:          e.Title,
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
		LabID:          e.LabID,
		Title:          e.Title,
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
