package lab_services

import (
	"context"
	"main/internal/dto"
	postgres_error "main/internal/infra/postgres"
	repository "main/internal/repository/lab"
	"main/internal/services/mapper"
)

type ChapterService interface {
	CreateChapter(ctx context.Context, chapter *dto.CreateChapterRequest, labID string) (*dto.ChapterResponse, error)
	GetChaptersByLabID(ctx context.Context, labID string) ([]dto.ChapterResponse, error)
	UpdateChapter(ctx context.Context, chapter *dto.UpdateChapterRequest) (*dto.ChapterResponse, error)
	DeleteChapter(ctx context.Context, chapterID string) error
}
type chapterService struct {
	chapterRepo repository.ChapterRepository
}

func NewChapterService(chapterRepo repository.ChapterRepository) ChapterService {
	return &chapterService{chapterRepo: chapterRepo}
}

func (c *chapterService) CreateChapter(ctx context.Context, req *dto.CreateChapterRequest, labID string) (*dto.ChapterResponse, error) {
	chapter := mapper.ToChapterModel(req, labID)
	res, err := c.chapterRepo.CreateChapter(ctx, chapter)
	if err != nil {
		return nil, postgres_error.MapError(err, "create chapter", "chapter")
	}
	return mapper.ToChapterResponse(res), nil
}

func (c *chapterService) DeleteChapter(ctx context.Context, chapterID string) error {
	err := c.chapterRepo.DeleteChapter(ctx, chapterID)
	if err != nil {
		return postgres_error.MapError(err, "delete chapter", "chapter")
	}
	return nil
}

func (c *chapterService) GetChaptersByLabID(ctx context.Context, labID string) ([]dto.ChapterResponse, error) {
	chapters, err := c.chapterRepo.GetChaptersByLabID(ctx, labID)
	if err != nil {
		return nil, postgres_error.MapError(err, "get chapters by lab ID", "chapter")
	}
	return mapper.ToChapterResponses(chapters), nil
}

func (c *chapterService) UpdateChapter(ctx context.Context, req *dto.UpdateChapterRequest) (*dto.ChapterResponse, error) {
	chapter := mapper.ToChapterModelFromUpdateRequest(req)
	chapter, err := c.chapterRepo.UpdateChapter(ctx, chapter)
	if err != nil {
		return nil, postgres_error.MapError(err, "update chapter", "chapter")
	}
	return mapper.ToChapterResponse(chapter), nil
}
