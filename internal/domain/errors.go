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
)

var AllErrorCodes = map[string]ErrorCode{
	string(CodeInternal):             CodeInternal,
	string(CodeInvalidRequest):       CodeInvalidRequest,
	string(CodeUnauthorized):         CodeUnauthorized,
	string(CodeForbidden):            CodeForbidden,
	string(CodeNotFound):             CodeNotFound,
	string(CodeConflict):             CodeConflict,
	string(CodeValidation):           CodeValidation,
	string(CodeUserNotFound):         CodeUserNotFound,
	string(CodeUserAlreadyExists):    CodeUserAlreadyExists,
	string(CodeSandboxNotFound):      CodeSandboxNotFound,
	string(CodeSandboxAlreadyExists): CodeSandboxAlreadyExists,
}
