package postgres_error

import (
	"fmt"
	"net/http"
	"strings"

	"main/internal/domain"
)

func MapError(err error, operation, entity string) *domain.AppError {
	if err == nil {
		return nil
	}

	switch {
	case IsRecordNotFound(err):
		return domain.NewAppError(
			http.StatusNotFound,
			entityCode(entity, "NOT_FOUND"),
			fmt.Sprintf("%s not found", entity),
			err,
			nil,
		)
	case IsUniqueViolation(err):
		return domain.NewAppError(
			http.StatusConflict,
			entityCode(entity, "ALREADY_EXISTS"),
			fmt.Sprintf("%s already exists", entity),
			err,
			nil,
		)
	case IsForeignKeyViolation(err):
		return domain.InvalidRequestError(
			fmt.Sprintf("%s reference invalid", entity),
			err,
		)
	case IsNotNullViolation(err):
		return domain.InvalidRequestError(
			fmt.Sprintf("%s missing required field", entity),
			err,
		)
	case IsCheckViolation(err):
		return domain.InvalidRequestError(
			fmt.Sprintf("%s check constraint violated", entity),
			err,
		)
	case IsOptimisticLock(err):
		return domain.NewAppError(
			http.StatusConflict,
			domain.CodeConflict,
			fmt.Sprintf("%s update conflict", entity),
			err,
			nil,
		)

	// Retryable / temporary failures
	case IsRetryable(err), IsConnectionError(err), IsQueryCanceled(err):
		return domain.NewAppError(
			http.StatusServiceUnavailable,
			domain.CodeInternal,
			fmt.Sprintf("Temporary %s failure, retry", operation),
			err,
			nil,
		)
	default:
		return domain.InternalError(err)
	}
}

func entityCode(entity, suffix string) domain.ErrorCode {
	codeStr := fmt.Sprintf("%s_%s", strings.ToUpper(entity), suffix)
	if c, ok := domain.AllErrorCodes[codeStr]; ok {
		return c
	}
	return domain.CodeInternal
}
