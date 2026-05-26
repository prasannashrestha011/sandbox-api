package mapper

import (
	"main/internal/dto"
	"main/internal/repository/model"
)

// UserCreateToDB maps a service create payload into a persistence model.
func UserCreateToDB(user *dto.UserCreate) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		Fullname: user.Fullname,
		Username: user.Username,
		Password: user.PasswordHash,
	}
}

// UserFromDB maps a persistence model into a service user.
func UserFromDB(user *model.User) *dto.User {
	if user == nil {
		return nil
	}

	return &dto.User{
		UserID:    user.UserID,
		Fullname:  user.Fullname,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UserUpdateToDB maps a service update payload into a persistence model.
func UserUpdateToDB(user *dto.UpdateUserInput) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		UserID:   user.UserID,
		Fullname: user.Fullname,
		Username: user.Username,
	}
}

// UserPasswordUpdateToDB maps a password update payload into a persistence model.
func UserPasswordUpdateToDB(update *dto.UserPasswordUpdate) *model.User {
	if update == nil {
		return nil
	}

	return &model.User{
		UserID:   update.UserID,
		Password: update.PasswordHash,
	}
}
