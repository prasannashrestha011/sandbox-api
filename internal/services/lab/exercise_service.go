package lab_services

import (
	"context"
	"main/internal/dto"
	postgres_error "main/internal/infra/postgres"
	repository "main/internal/repository/lab"
	"main/internal/services/mapper"
)

type ExerciseService interface {
	CreateExercise(ctx context.Context, exercise *dto.CreateExerciseRequest, chapterID string) (*dto.ExerciseResponse, error)
	GetExerciseByID(ctx context.Context, id string) (*dto.ExerciseResponse, error)
	UpdateExercise(ctx context.Context, id string, exercise *dto.UpdateExerciseRequest) error
	DeleteExercise(ctx context.Context, id string) error
	ListExercisesByChapterID(ctx context.Context, chapterID string) ([]*dto.ExerciseResponse, error)
}

type exerciseService struct {
	exerciseRepo repository.ExerciseRepository
}

// CreateExercise implements [ExerciseService].
func (e *exerciseService) CreateExercise(ctx context.Context, exercise *dto.CreateExerciseRequest, chapterID string) (*dto.ExerciseResponse, error) {
	exerciseModel := mapper.ToExerciseModel(exercise, chapterID)
	resp, err := e.exerciseRepo.CreateExercise(ctx, exerciseModel)
	if err != nil {
		return nil, postgres_error.MapError(err, "create exercise", "exercise")
	}

	return mapper.ToExerciseResponse(resp), nil
}

// DeleteExercise implements [ExerciseService].
func (e *exerciseService) DeleteExercise(ctx context.Context, id string) error {
	err := e.exerciseRepo.DeleteExercise(ctx, id)
	if err != nil {
		return postgres_error.MapError(err, "delete exercise", "exercise")
	}
	return nil
}

// GetExerciseByID implements [ExerciseService].
func (e *exerciseService) GetExerciseByID(ctx context.Context, id string) (*dto.ExerciseResponse, error) {
	exerciseModel, err := e.exerciseRepo.GetExerciseByID(ctx, id)
	if err != nil {
		return nil, postgres_error.MapError(err, "get exercise", "exercise")
	}
	return mapper.ToExerciseResponse(exerciseModel), nil
}

// ListExercisesByChapterID implements [ExerciseService].
func (e *exerciseService) ListExercisesByChapterID(ctx context.Context, chapterID string) ([]*dto.ExerciseResponse, error) {
	exercisesModels, err := e.exerciseRepo.ListExercisesByChapterID(ctx, chapterID)
	if err != nil {
		return nil, postgres_error.MapError(err, "list exercises", "exercise")
	}
	resp := make([]*dto.ExerciseResponse, len(exercisesModels))
	for i, model := range exercisesModels {
		resp[i] = mapper.ToExerciseResponse(model)
	}
	return resp, nil
}

// UpdateExercise implements [ExerciseService].
func (e *exerciseService) UpdateExercise(ctx context.Context, id string, exercise *dto.UpdateExerciseRequest) error {
	exerciseModel := mapper.ToExerciseModelFromUpdateRequest(exercise, id)
	err := e.exerciseRepo.UpdateExercise(ctx, exerciseModel)
	if err != nil {
		return postgres_error.MapError(err, "update exercise", "exercise")
	}
	return nil
}

func NewExerciseService(exerciseRepo repository.ExerciseRepository) ExerciseService {
	return &exerciseService{
		exerciseRepo: exerciseRepo,
	}
}
