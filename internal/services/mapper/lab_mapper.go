package mapper

import (
	"context"
	request_context "main/internal/context"
	"main/internal/dto"
	"main/internal/services/models"
)

// ToLabModel maps a DTO create request to a service Lab model.
func ToLabModel(req *dto.CreateLabRequest, ctx context.Context) *models.Lab {
	if req == nil {
		return nil
	}
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil
	}

	exercises := make([]models.Exercise, len(req.Exercises))
	for i, ex := range req.Exercises {
		exercises[i] = *ToExerciseModel(&ex)
	}
	tags := make([]models.Tag, len(req.Tags))
	for i, t := range req.Tags {
		tags[i] = models.Tag{Name: t}
	}

	return &models.Lab{
		Title:       req.Title,
		Description: req.Description,
		Lang:        req.Lang,
		Difficulty:  req.Difficulty,
		IsPublic:    req.IsPublic,
		Tags:        tags,
		Exercises:   exercises,
		ContainerID: req.ContainerID,
		CreatedByID: userID.String(),
	}
}

// ToLabResponse maps a service Lab model to a DTO Lab response.
func ToLabResponse(l *models.Lab) *dto.LabResponse {
	if l == nil {
		return nil
	}

	exercises := make([]dto.ExerciseResponse, len(l.Exercises))
	for i, ex := range l.Exercises {
		exercises[i] = *ToExerciseResponse(&ex)
	}

	tags := make([]string, len(l.Tags))
	for i, t := range l.Tags {
		tags[i] = t.Name
	}

	return &dto.LabResponse{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		Lang:        l.Lang,
		Difficulty:  l.Difficulty,
		IsPublic:    l.IsPublic,
		ContainerID: l.ContainerID,
		CreatedByID: l.CreatedByID,
		Tags:        tags,
		Exercises:   exercises,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

// ToExerciseModel maps a DTO create exercise request to a service Exercise model.
func ToExerciseModel(req *dto.CreateExerciseRequest) *models.Exercise {
	if req == nil {
		return nil
	}
	return &models.Exercise{
		Title:          req.Title,
		Description:    req.Description,
		StarterCode:    req.StarterCode,
		ExpectedOutput: req.ExpectedOutput,
		Hints:          req.Hints,
		OrderIndex:     req.OrderIndex,
		Solution:       req.Solution,
		MaxAttempts:    req.MaxAttempts,
	}
}

// ToExerciseResponse maps a service Exercise model to a DTO Exercise response.
func ToExerciseResponse(e *models.Exercise) *dto.ExerciseResponse {
	if e == nil {
		return nil
	}
	return &dto.ExerciseResponse{
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
