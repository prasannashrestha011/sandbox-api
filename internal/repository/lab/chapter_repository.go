package repository

import (
	"context"
	"main/internal/repository/mapper"
	lab_model "main/internal/repository/model/lab"
	"main/internal/services/models"

	"gorm.io/gorm"
)

type chapterRepository struct {
	db *gorm.DB
}

type ChapterRepository interface {
	CreateChapter(ctx context.Context, req *models.Chapter) (*models.Chapter, error)
	GetChaptersByLabID(ctx context.Context, labID string) ([]models.Chapter, error)
	UpdateChapter(ctx context.Context, chapter *models.Chapter) (*models.Chapter, error)
	DeleteChapter(ctx context.Context, chapterID string) error
}

// CreateChapter implements [ChapterRepository].
func (c *chapterRepository) CreateChapter(ctx context.Context, req *models.Chapter) (*models.Chapter, error) {
	chapter := mapper.ChapterToGorm(req)
	if err := c.db.WithContext(ctx).Model(&lab_model.Chapter{}).Create(chapter).Error; err != nil {
		return nil, err
	}
	res := mapper.ChapterFromGorm(chapter)
	return res, nil
}

// DeleteChapter implements [ChapterRepository].
func (c *chapterRepository) DeleteChapter(ctx context.Context, chapterID string) error {
	if err := c.db.WithContext(ctx).Model(&lab_model.Chapter{}).Delete(&lab_model.Chapter{}, "id = ?", chapterID).Error; err != nil {
		return err
	}
	return nil
}

// GetChaptersByLabID implements [ChapterRepository].
func (c *chapterRepository) GetChaptersByLabID(ctx context.Context, labID string) ([]models.Chapter, error) {
	var chapters []models.Chapter
	if err := c.db.WithContext(ctx).Model(&lab_model.Chapter{}).Where("lab_id = ?", labID).Find(&chapters).Error; err != nil {
		return nil, err
	}
	return chapters, nil
}

// UpdateChapter implements [ChapterRepository].
func (c *chapterRepository) UpdateChapter(ctx context.Context, chapter *models.Chapter) (*models.Chapter, error) {
	if err := c.db.WithContext(ctx).Model(&lab_model.Chapter{}).Save(chapter).Error; err != nil {
		return nil, err
	}
	return chapter, nil
}

func NewChapterRepository(db *gorm.DB) ChapterRepository {
	return &chapterRepository{db: db}
}
