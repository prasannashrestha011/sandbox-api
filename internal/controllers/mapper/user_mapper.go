package mapper

import (
	"main/internal/dto"
	"main/internal/enums"
	"main/internal/repository/model"
)

// UserModelToDTO maps a domain user to an API-safe DTO.
func UserModelToDTO(user *model.User) *dto.User {
	if user == nil {
		return nil
	}

	return &dto.User{
		UserID:    user.UserID,
		Fullname:  user.Fullname,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UserModelsToDTO maps a slice of domain users to DTOs.
func UserModelsToDTO(users []model.User) []dto.User {
	items := make([]dto.User, 0, len(users))
	for i := range users {
		mapped := UserModelToDTO(&users[i])
		if mapped != nil {
			items = append(items, *mapped)
		}
	}
	return items
}

// CreateUserInputToModel maps a create request DTO to the domain model.
func CreateUserInputToModel(input *dto.CreateUserInput) *model.User {
	if input == nil {
		return nil
	}

	return &model.User{
		Fullname: input.Fullname,
		Username: input.Username,
		Password: input.Password,
		Role:     enums.Role(input.Role),
	}
}

// UpdateUserInputToModel maps an update request DTO to the domain model.
func UpdateUserInputToModel(input *dto.UpdateUserInput) *model.User {
	if input == nil {
		return nil
	}

	return &model.User{
		UserID:   input.UserID,
		Fullname: input.Fullname,
		Username: input.Username,
		Role:     enums.Role(input.Role),
	}
}
