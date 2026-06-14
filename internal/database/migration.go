package database

import (
	"main/internal/repository/model"
	lab_model "main/internal/repository/model/lab"

	"gorm.io/gorm"
)

// AutoMigrate runs the schema migrations for persistence models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.SandboxTemplate{}, &model.SandboxSession{}, &model.DockerImage{}, &model.User{}, &model.RefreshToken{},
		&lab_model.Lab{}, &lab_model.Exercise{}, &lab_model.Chapter{}, &lab_model.Exercise{}, &lab_model.LabEnrollment{}, &lab_model.Submission{}, &lab_model.Tag{})
}
