package lab_services

import (
	"context"
	"fmt"
	"log"
	request_context "main/internal/context"
	"main/internal/domain"
	"main/internal/dto"
	"main/internal/enums"
	postgres_error "main/internal/infra/postgres"
	repository "main/internal/repository/lab"

	"main/internal/services/mapper"
)

type EnrollmentService interface {
	EnrollUserToLab(ctx context.Context, req *dto.EnrollmentRequest) error
	GetEnrollment(ctx context.Context, userID, labID string) (*dto.EnrollmentResponse, error)
	GetUserEnrollments(ctx context.Context, userID string) ([]dto.EnrollmentResponse, error)
	DeleteEnrollment(ctx context.Context, userID, labID string) error
}

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository) EnrollmentService {
	return &enrollmentService{enrollmentRepo: enrollmentRepo}
}

// EnrollUserToLab implements [EnrollmentService].
func (e *enrollmentService) EnrollUserToLab(ctx context.Context, req *dto.EnrollmentRequest) error {
	userType, ok := request_context.UserType(ctx)
	log.Println("User Type", userType)
	if !ok {
		return postgres_error.MapError(fmt.Errorf("user type not found"), "enroll user", "enrollment")
	}
	if userType != enums.UserTypeStudent {
		return domain.InvalidRequestError("only students can enroll to labs", nil)
	}
	enrollment := mapper.ToEnrollmentToModel(ctx, req)
	err := e.enrollmentRepo.EnrollUserToLab(ctx, enrollment)
	if err != nil {
		return postgres_error.MapError(err, "enroll user", "enrollment")
	}
	return nil
}

// GetEnrollment implements [EnrollmentService].
func (e *enrollmentService) GetEnrollment(ctx context.Context, userID string, labID string) (*dto.EnrollmentResponse, error) {
	enrollment, err := e.enrollmentRepo.GetEnrollment(ctx, userID, labID)
	if err != nil {
		return nil, postgres_error.MapError(err, "get enrollment", "enrollment")
	}
	return mapper.ToEnrollmentResponse(enrollment), nil
}

// GetUserEnrollments implements [EnrollmentService].
func (e *enrollmentService) GetUserEnrollments(ctx context.Context, userID string) ([]dto.EnrollmentResponse, error) {
	enrollments, err := e.enrollmentRepo.GetUserEnrollments(ctx, userID)
	if err != nil {
		return nil, postgres_error.MapError(err, "get user enrollments", "enrollment")
	}
	return mapper.ToEnrollmentResponses(enrollments), nil
}

// DeleteEnrollment implements [EnrollmentService].
func (e *enrollmentService) DeleteEnrollment(ctx context.Context, userID string, labID string) error {
	err := e.enrollmentRepo.DeleteEnrollment(ctx, userID, labID)
	if err != nil {
		return postgres_error.MapError(err, "delete enrollment", "enrollment")
	}
	return nil
}
