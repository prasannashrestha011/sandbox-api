package dto

import "time"

type CreateLabRequest struct {
	Title       string                  `json:"title" binding:"required"`
	Description string                  `json:"description" binding:"required"`
	Lang        string                  `json:"lang" binding:"required"`
	Difficulty  string                  `json:"difficulty" binding:"required"`
	IsPublic    bool                    `json:"isPublic"`
	Exercises   []CreateExerciseRequest `json:"exercises"`
	Tags        []string                `json:"tags"`
}

type UpdateLabRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Lang        *string  `json:"lang,omitempty"`
	Difficulty  *string  `json:"difficulty,omitempty"`
	IsPublic    *bool    `json:"isPublic,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type LabResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Lang        string             `json:"lang"`
	Difficulty  string             `json:"difficulty"`
	IsPublic    bool               `json:"isPublic"`
	ContainerID string             `json:"containerId"`
	CreatedByID string             `json:"createdById"`
	Tags        []string           `json:"tags"`
	Exercises   []ExerciseResponse `json:"exercises,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type CreateExerciseRequest struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description" binding:"required"`
	StarterCode    string `json:"starterCode"`
	ExpectedOutput string `json:"expectedOutput" binding:"required"`
	Hints          string `json:"hints"`
	OrderIndex     int    `json:"orderIndex"`
	Solution       string `json:"solution"`
	MaxAttempts    int    `json:"maxAttempts"`
}

type UpdateExerciseRequest struct {
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	StarterCode    *string `json:"starterCode,omitempty"`
	ExpectedOutput *string `json:"expectedOutput,omitempty"`
	Hints          *string `json:"hints,omitempty"`
	OrderIndex     *int    `json:"orderIndex,omitempty"`
	Solution       *string `json:"solution,omitempty"`
	MaxAttempts    *int    `json:"maxAttempts,omitempty"`
}

type ExerciseResponse struct {
	ID             string    `json:"id"`
	LabID          string    `json:"labId"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	StarterCode    string    `json:"starterCode"`
	ExpectedOutput string    `json:"expectedOutput"`
	Hints          string    `json:"hints"`
	OrderIndex     int       `json:"orderIndex"`
	Solution       string    `json:"solution,omitempty"`
	MaxAttempts    int       `json:"maxAttempts"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
