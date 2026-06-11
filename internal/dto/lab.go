package dto

import (
	"strings"
	"time"
)

type CreateLabRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Lang        string   `json:"lang" binding:"required"`
	Difficulty  string   `json:"difficulty" binding:"required"`
	IsPublic    bool     `json:"isPublic"`
	Tags        []string `json:"tags"`
	ContainerID string   `json:"container_id" binding:"required"`
}

func (c *CreateLabRequest) Sanitize() {
	c.Title = strings.TrimSpace(c.Title)
	c.Description = strings.TrimSpace(c.Description)
	c.Lang = strings.ToLower(strings.TrimSpace(c.Lang))
	c.Difficulty = strings.ToLower(strings.TrimSpace(c.Difficulty))
	c.ContainerID = strings.TrimSpace(c.ContainerID)
}
func (c *CreateLabRequest) Validate() error {
	var v ValidationErrors
	c.Sanitize()

	if c.Title == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "title",
			Message: "title is required",
		})
	}
	if c.Description == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "description",
			Message: "description is required",
		})
	}
	if c.Lang == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "lang",
			Message: "lang is required",
		})
	}
	if c.Difficulty == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "difficulty",
			Message: "difficulty is required",
		})
	}
	if c.ContainerID == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "containerID",
			Message: "containerID is required",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

type UpdateLabRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Lang        *string  `json:"lang,omitempty"`
	Difficulty  *string  `json:"difficulty,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type LabResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Lang        string            `json:"lang"`
	Difficulty  string            `json:"difficulty"`
	IsPublic    bool              `json:"is_public"`
	ContainerID string            `json:"container_id"`
	CreatedByID string            `json:"created_by_id"`
	Tags        []string          `json:"tags"`
	Chapters    []ChapterResponse `json:"chapters,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// ##
type CreateChapterRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	OrderIndex  int    `json:"order_index"`
}
type CreateChapterResponse struct {
	ID string `json:"id"`
}
type ChapterResponse struct {
	ID          string             `json:"id"`
	LabID       string             `json:"labId"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	OrderIndex  int                `json:"order_index"`
	Exercises   []ExerciseResponse `json:"exercises,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}
type UpdateChapterRequest struct {
	ID          string  `json:"id" binding:"required"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	OrderIndex  *int    `json:"order_index,omitempty"`
}

// ##
type CreateExerciseRequest struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description" binding:"required"`
	StarterCode    string `json:"starterCode"`
	ExpectedOutput string `json:"expectedOutput" binding:"required"`
	Hints          string `json:"hints"`
	OrderIndex     int    `json:"order_index"`
	Solution       string `json:"solution"`
	MaxAttempts    int    `json:"maxAttempts"`
}
type CreateExerciseResponse struct {
	ID string `json:"id"`
}

type UpdateExerciseRequest struct {
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	StarterCode    *string `json:"starterCode,omitempty"`
	ExpectedOutput *string `json:"expectedOutput,omitempty"`
	Hints          *string `json:"hints,omitempty"`
	OrderIndex     *int    `json:"order_index,omitempty"`
	Solution       *string `json:"solution,omitempty"`
	MaxAttempts    *int    `json:"maxAttempts,omitempty"`
}

type ExerciseResponse struct {
	ID             string    `json:"id"`
	ChapterID      string    `json:"chapterID"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	StarterCode    string    `json:"starterCode"`
	ExpectedOutput string    `json:"expectedOutput"`
	Hints          string    `json:"hints"`
	OrderIndex     int       `json:"order_index"`
	Solution       string    `json:"solution,omitempty"`
	MaxAttempts    int       `json:"maxAttempts"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// ##

type EnrollmentRequest struct {
	LabID string `json:"labId" binding:"required"`
}
type EnrollmentResponse struct {
	UserID      string     `json:"userID" binding:"required"`
	LabID       string     `json:"labId" binding:"required"`
	Status      string     `json:"status"` // enrolled | in_progress | completed
	ProgressPct int        `json:"progressPct"`
	EnrolledAt  time.Time  `json:"enrolledAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

type SubmissionRequest struct {
	Code   string `json:"code" binding:"required"`
	Output string `json:"output,omitempty"`
}

func (s *SubmissionRequest) Sanitize() {
	s.Code = strings.TrimSpace(s.Code)
}
func (s *SubmissionRequest) Validate() error {
	var v ValidationErrors
	s.Sanitize()

	if s.Code == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "code",
			Message: "code is required",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

type SubmissionResponse struct {
	ID          string    `json:"id"`
	ExerciseID  string    `json:"exerciseId"`
	UserID      string    `json:"userId"`
	Language    string    `json:"language"`
	Code        string    `json:"code"`
	Output      string    `json:"output,omitempty"`
	Status      string    `json:"status"` // pending | running | completed
	Score       int       `json:"score,omitempty"`
	AttemptNo   int       `json:"attemptNo"`
	SubmittedAt time.Time `json:"submittedAt"`
}
