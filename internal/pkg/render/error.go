package render

import (
	"fmt"
	"strings"
)

const (
	// InternalServerError ...
	InternalServerError string = "internal_server_error"
	// NotFoundError ...
	NotFoundError string = "not_found_error"
	// BadRequestError ...
	BadRequestError string = "bad_request_error"
	// UnauthorizedError ...
	UnauthorizedError string = "unauthorized_error"
	// RequestValidationError ...
	RequestValidationError string = "validation_error"
	// ConflictError ...
	ConflictError string = "conflict_error"
	// MethodNotAllowedError ...
	MethodNotAllowedError string = "method_not_allowed_error"
)

// Error is a generic error response.
type Error struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Errors  *FieldErrors `json:"errors"`
}

// FieldErrors is a list of field errors.
type FieldErrors []*FieldError

// Error is the implementation of the error interface.
func (e *FieldErrors) Error() string {
	errMsg := &strings.Builder{}
	for _, fieldErr := range *e {
		errMsg.WriteString(fieldErr.Error()) //nolint: revive //WriteString can't return a non-nil error
	}

	return errMsg.String()
}

// Add appends a new error to the list of field errors.
func (e *FieldErrors) Add(field, message string, err error) {
	*e = append(*e, &FieldError{
		Field:   field,
		Message: message,
		err:     err,
	})
}

// FieldError represents a single field error, optinally wrapping another error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	err     error
}

// Error is the implementation of the error interface.
func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Unwrap returns the wrapped error (if present).
func (e *FieldError) Unwrap() error {
	return e.err
}

// NewError creates a new error response.
func NewError(status, message string, errs ...*FieldError) *Error {
	return &Error{
		Status:  status,
		Message: message,
		Errors:  (*FieldErrors)(&errs),
	}
}

// NewValidationError creates a new validation error response.
func NewValidationError(errs ...*FieldError) *Error {
	errsPtr := FieldErrors(errs)
	return &Error{
		Status:  RequestValidationError,
		Message: "validation error",
		Errors:  &errsPtr,
	}
}

// Error is the implementation of the error interface.
func (e *Error) Error() string {
	return e.Message
}
