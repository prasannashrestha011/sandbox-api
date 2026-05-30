package postgres_error

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// GORM level helpers
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsOptimisticLock(err error) bool {
	return false
}

// PostgreSQL error helpers (using GORM's translated errors)
func IsUniqueViolation(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}

func IsForeignKeyViolation(err error) bool {
	return errors.Is(err, gorm.ErrForeignKeyViolated)
}

func IsNotNullViolation(err error) bool {
	return errors.Is(err, gorm.ErrInvalidData)
}

func IsCheckViolation(err error) bool {
	return errors.Is(err, gorm.ErrCheckConstraintViolated)
}

// Concurrency or retryable errors
func IsSerializationFailure(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "serialization failure") ||
		strings.Contains(errMsg, "SQLSTATE 40001")
}

func IsDeadlock(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "deadlock") ||
		strings.Contains(errMsg, "SQLSTATE 40P01")
}

func IsRetryable(err error) bool {
	return IsSerializationFailure(err) || IsDeadlock(err)
}

// Infrastructure failures
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "SQLSTATE 08")
}

func IsQueryCanceled(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "canceled") ||
		strings.Contains(errMsg, "SQLSTATE 57014")
}
