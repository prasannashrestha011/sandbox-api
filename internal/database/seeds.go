package database

import (
	lab_model "main/internal/repository/model/lab"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&lab_model.Tag{}).Count(&count)

	if count > 0 {
		return nil
	}

	tags := []lab_model.Tag{
		{Name: "Go"},
		{Name: "Python"},
		{Name: "JavaScript"},
	}

	return db.Create(&tags).Error

}
