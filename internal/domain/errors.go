package domain

import "errors"

var ErrDatabaseError = errors.New("database error:")

type ErrorCode string

const (
	// Generic
	CodeInternal       ErrorCode = "INTERNAL_ERROR"
	CodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	CodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	CodeForbidden      ErrorCode = "FORBIDDEN"
	CodeNotFound       ErrorCode = "NOT_FOUND"
	CodeConflict       ErrorCode = "CONFLICT"
	CodeValidation     ErrorCode = "VALIDATION_ERROR"

	//User
	CodeUserNotFound      ErrorCode = "USER_NOT_FOUND"
	CodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"

	//Sandbox
	CodeSandboxNotFound      ErrorCode = "SANDBOX_NOT_FOUND"
	CodeSandboxAlreadyExists ErrorCode = "SANDBOX_ALREADY_EXISTS"

	//DockerImage
	CodeDockerImageNotFound      ErrorCode = "DOCKER_IMAGE_NOT_FOUND"
	CodeDockerImageAlreadyExists ErrorCode = "DOCKER_IMAGE_ALREADY_EXISTS"
	//Lab
	CodeLabNotFound      ErrorCode = "LAB_NOT_FOUND"
	CodeLabAlreadyExists ErrorCode = "LAB_ALREADY_EXISTS"

	CodeChapterNotFound      ErrorCode = "CHAPTER_NOT_FOUND"
	CodeChapterAlreadyExists ErrorCode = "CHAPTER_ALREADY_EXISTS"

	CodeExerciseNotFound      ErrorCode = "EXERCISE_NOT_FOUND"
	CodeExerciseAlreadyExists ErrorCode = "EXERCISE_ALREADY_EXISTS"
)

var AllErrorCodes = map[string]ErrorCode{
	string(CodeInternal):                 CodeInternal,
	string(CodeInvalidRequest):           CodeInvalidRequest,
	string(CodeUnauthorized):             CodeUnauthorized,
	string(CodeForbidden):                CodeForbidden,
	string(CodeNotFound):                 CodeNotFound,
	string(CodeConflict):                 CodeConflict,
	string(CodeValidation):               CodeValidation,
	string(CodeUserNotFound):             CodeUserNotFound,
	string(CodeUserAlreadyExists):        CodeUserAlreadyExists,
	string(CodeSandboxNotFound):          CodeSandboxNotFound,
	string(CodeSandboxAlreadyExists):     CodeSandboxAlreadyExists,
	string(CodeDockerImageNotFound):      CodeDockerImageNotFound,
	string(CodeDockerImageAlreadyExists): CodeDockerImageAlreadyExists,
	string(CodeLabNotFound):              CodeLabNotFound,
	string(CodeLabAlreadyExists):         CodeLabAlreadyExists,
	string(CodeChapterNotFound):          CodeChapterNotFound,
	string(CodeChapterAlreadyExists):     CodeChapterAlreadyExists,
	string(CodeExerciseNotFound):         CodeExerciseNotFound,
	string(CodeExerciseAlreadyExists):    CodeExerciseAlreadyExists,
}
