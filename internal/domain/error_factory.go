package domain

import "net/http"

// Generic errors and helpers
func InternalError(err error) *AppError {
	return NewAppError(
		http.StatusInternalServerError,
		CodeInternal,
		"something went wrong",
		err,
		nil,
	)
}

func InvalidRequestError(msg string, err error) *AppError {
	return NewAppError(
		http.StatusBadRequest,
		CodeInvalidRequest,
		msg,
		err,
		nil,
	)
}

func UnauthorizedError(msg string) *AppError {
	return NewAppError(
		http.StatusUnauthorized,
		CodeUnauthorized,
		msg,
		nil,
		nil,
	)
}

func ForbiddenError(err error) *AppError {
	return NewAppError(
		http.StatusForbidden,
		CodeForbidden,
		"forbidden",
		err,
		nil,
	)
}

func NotFoundError(err error) *AppError {
	return NewAppError(
		http.StatusNotFound,
		CodeNotFound,
		"not found",
		err,
		nil,
	)
}

func ConflictError(err error) *AppError {
	return NewAppError(
		http.StatusConflict,
		CodeConflict,
		"conflict",
		err,
		nil,
	)
}

func ValidationError(err error, errDetails any) *AppError {
	return NewAppError(
		http.StatusUnprocessableEntity,
		CodeValidation,
		"validation failed",
		err,
		errDetails,
	)
}

// helper method for creating a new AppError
func NewAppError(
	httpCode int,
	code ErrorCode,
	message string,
	err error,
	errorDetails any,
) *AppError {
	return &AppError{
		HTTPCode:     httpCode,
		Code:         code,
		Message:      message,
		Err:          err,
		ErrorDetails: errorDetails,
	}
}
