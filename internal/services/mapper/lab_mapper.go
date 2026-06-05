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
		ContainerID: req.ContainerID,
		CreatedByID: userID.String(),
	}
}

// ToLabResponse maps a service Lab model to a DTO Lab response.
func ToLabResponse(l *models.Lab) *dto.LabResponse {
	if l == nil {
		return nil
	}

	tags := make([]string, len(l.Tags))
	for i, t := range l.Tags {
		tags[i] = t.Name
	}
	chapters := make([]dto.ChapterResponse, len(l.Chapters))
	for i, c := range l.Chapters {
		chapters[i] = *ToChapterResponse(&c)
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
		Chapters:    chapters,

		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}

// ToExerciseModel maps a DTO create exercise request to a service Exercise model.
func ToExerciseModel(req *dto.CreateExerciseRequest, chapterID string) *models.Exercise {
	if req == nil {
		return nil
	}
	return &models.Exercise{
		ChapterID:      chapterID,
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
func ToExerciseModelFromUpdateRequest(req *dto.UpdateExerciseRequest, id string) *models.Exercise {
	if req == nil {
		return nil
	}
	ex := &models.Exercise{
		ID: id,
	}
	if req.Title != nil {
		ex.Title = *req.Title
	}
	if req.Description != nil {
		ex.Description = *req.Description
	}
	if req.StarterCode != nil {
		ex.StarterCode = *req.StarterCode
	}
	if req.ExpectedOutput != nil {
		ex.ExpectedOutput = *req.ExpectedOutput
	}
	if req.Hints != nil {
		ex.Hints = *req.Hints
	}
	if req.OrderIndex != nil {
		ex.OrderIndex = *req.OrderIndex
	}
	if req.Solution != nil {
		ex.Solution = *req.Solution
	}
	if req.MaxAttempts != nil {
		ex.MaxAttempts = *req.MaxAttempts
	}

	return ex
}

// ToExerciseResponse maps a service Exercise model to a DTO Exercise response.
func ToExerciseResponse(e *models.Exercise) *dto.ExerciseResponse {
	if e == nil {
		return nil
	}
	return &dto.ExerciseResponse{
		ID:             e.ID,
		ChapterID:      e.ChapterID,
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

func ToChapterModel(req *dto.CreateChapterRequest, labID string) *models.Chapter {
	if req == nil {
		return nil
	}

	return &models.Chapter{
		Title:       req.Title,
		Description: req.Description,
		OrderIndex:  req.OrderIndex,
		LabID:       labID,
	}
}

func ToChapterResponse(c *models.Chapter) *dto.ChapterResponse {
	if c == nil {
		return nil
	}

	exercises := make([]dto.ExerciseResponse, len(c.Exercises))
	for i, e := range c.Exercises {
		exercises[i] = *ToExerciseResponse(&e)
	}

	return &dto.ChapterResponse{
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
func ToChapterResponses(chapters []models.Chapter) []dto.ChapterResponse {
	res := make([]dto.ChapterResponse, len(chapters))
	for i, c := range chapters {
		res[i] = *ToChapterResponse(&c)
	}
	return res
}

func ToChapterModelFromUpdateRequest(req *dto.UpdateChapterRequest) *models.Chapter {
	if req == nil {
		return nil
	}

	ch := &models.Chapter{
		ID: req.ID,
	}

	if req.Title != nil {
		ch.Title = *req.Title
	}
	if req.Description != nil {
		ch.Description = *req.Description
	}
	if req.OrderIndex != nil {
		ch.OrderIndex = *req.OrderIndex
	}

	return ch
}

func ToEnrollmentToModel(ctx context.Context, req *dto.EnrollmentRequest) *models.LabEnrollment {

	if req == nil {
		return nil
	}
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil
	}
	return &models.LabEnrollment{
		UserID: userID.String(),
		LabID:  req.LabID,
	}
}

func ToEnrollmentResponse(e *models.LabEnrollment) *dto.EnrollmentResponse {
	if e == nil {
		return nil
	}
	return &dto.EnrollmentResponse{
		UserID:      e.UserID,
		LabID:       e.LabID,
		Status:      e.Status,
		ProgressPct: e.ProgressPct,
		EnrolledAt:  e.EnrolledAt,
		CompletedAt: e.CompletedAt,
	}
}

func ToEnrollmentResponses(enrollments []models.LabEnrollment) []dto.EnrollmentResponse {
	res := make([]dto.EnrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		res[i] = *ToEnrollmentResponse(&e)
	}
	return res
}

func ToSubmissionModel(ctx context.Context, req *dto.SubmissionRequest) *models.Submission {
	if req == nil {
		return nil
	}
	userID, ok := request_context.UserID(ctx)
	if !ok {
		return nil
	}
	return &models.Submission{
		UserID:     userID.String(),
		ExerciseID: req.ExerciseID,
		Code:       req.Code,
	}
}
func ToSubmissionResponse(s *models.Submission) *dto.SubmissionResponse {
	if s == nil {
		return nil
	}
	return &dto.SubmissionResponse{
		ID:          s.ID,
		UserID:      s.UserID,
		ExerciseID:  s.ExerciseID,
		Code:        s.Code,
		Language:    s.Language,
		Output:      s.Output,
		Status:      s.Status,
		Score:       s.Score,
		AttemptNo:   s.AttemptNo,
		SubmittedAt: s.SubmittedAt,
	}
}

func ToSubmissionResponses(submissions []models.Submission) []dto.SubmissionResponse {
	res := make([]dto.SubmissionResponse, len(submissions))
	for i, s := range submissions {
		res[i] = *ToSubmissionResponse(&s)
	}
	return res
}
