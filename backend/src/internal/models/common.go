package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type ID uint64

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

func GetCurrentDate() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

func NewValidationError(field string, message string) ValidationError {
	return ValidationError{field, message}
}

func NewValidationErrorf(field string, format string, args ...interface{}) ValidationError {
	return ValidationError{field, fmt.Sprintf(format, args...)}
}

func IsValidationError(err error) bool {
	return IsCustomError[ValidationError](err)
}

func IsCustomError[E error](err error) bool {
	var e E
	ok := errors.As(err, &e)
	return ok
}
