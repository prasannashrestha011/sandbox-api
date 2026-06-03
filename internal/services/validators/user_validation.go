package services_validators

import (
	"main/internal/dto"
	"main/internal/enums"
	"main/internal/repository/model"
)

func ValidateUserTypeAndRole(req *model.User) error {
	var v dto.ValidationErrors
	if req.UserType == enums.UserTypeStudent && req.Role == enums.RoleAdmin {
		v.Violations = append(v.Violations, dto.FieldViolation{
			Field:   "role",
			Message: "students cannot have admin role",
		})
	}
	if len(v.Violations) > 0 {
		return &v
	}
	return nil

}
