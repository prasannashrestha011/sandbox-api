package repository

import (
	"main/internal/repository/mapper"
	gormodel "main/internal/repository/model"
	"main/internal/services/models"

	"gorm.io/gorm"
)

type WarmPoolRepository interface {
	CreateWarmPool(req *models.WarmPool) (*models.WarmPool, error)
	GetWarmPool(id string) (*models.WarmPool, error)
}
type warmPoolRepository struct {
	db *gorm.DB
}

func NewWarmPoolRepository(db *gorm.DB) WarmPoolRepository {
	return &warmPoolRepository{
		db: db,
	}
}

func (w *warmPoolRepository) CreateWarmPool(req *models.WarmPool) (*models.WarmPool, error) {

	warmPool := mapper.WarmPoolToGorm(req)
	err := w.db.Model(&gormodel.WarmPool{}).Create(warmPool).Error
	if err != nil {
		return nil, err
	}
	return mapper.WarmPoolFromGorm(warmPool), nil
}

// GetWarmPool implements [WarmPoolRepository].
func (w *warmPoolRepository) GetWarmPool(id string) (*models.WarmPool, error) {
	var warmPool gormodel.WarmPool
	err := w.db.Model(&gormodel.WarmPool{}).Where("id = ?", id).First(&warmPool).Error
	if err != nil {
		return nil, err
	}
	return mapper.WarmPoolFromGorm(&warmPool), nil
}
