// Package errmgt provides utilities for managing and handling errors in Go applications.
package errmgt

import (
	"errors"
	"fmt"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	// ValidationError represents input validation errors
	ValidationError ErrorType = "validation"
	// BusinessError represents business logic errors
	BusinessError ErrorType = "business"
	// SystemError represents system-level errors
	SystemError ErrorType = "system"
	// ExternalError represents errors from external services
	ExternalError ErrorType = "external"
)

// ManagedError is a structured error with additional context
type ManagedError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    string            `json:"details,omitempty"`
	Cause      error             `json:"-"`
	Context    map[string]string `json:"context,omitempty"`
	Type       ErrorType         `json:"type"`
	StatusCode int               `json:"status_code,omitempty"`
	Retryable  bool              `json:"retryable"`
}

// Error implements the error interface
func (e *ManagedError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s:%s] %s: %s", e.Type, e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *ManagedError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target error
func (e *ManagedError) Is(target error) bool {
	if target == nil {
		return false
	}

	var managedErr *ManagedError
	if errors.As(target, &managedErr) {
		return e.Type == managedErr.Type && e.Code == managedErr.Code
	}

	return errors.Is(e.Cause, target)
}

// NewError creates a new ManagedError
func NewError(errType ErrorType, code, message string) *ManagedError {
	return &ManagedError{
		Type:    errType,
		Code:    code,
		Message: message,
		Context: make(map[string]string),
	}
}

// NewErrorWithCause creates a new ManagedError wrapping an existing error
func NewErrorWithCause(errType ErrorType, code, message string, cause error) *ManagedError {
	return &ManagedError{
		Type:    errType,
		Code:    code,
		Message: message,
		Cause:   cause,
		Context: make(map[string]string),
	}
}

// WithDetails adds details to the error
func (e *ManagedError) WithDetails(details string) *ManagedError {
	e.Details = details
	return e
}

// WithContext adds context information to the error
func (e *ManagedError) WithContext(key, value string) *ManagedError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// WithRetryable sets whether the error is retryable
func (e *ManagedError) WithRetryable(retryable bool) *ManagedError {
	e.Retryable = retryable
	return e
}

// WithStatusCode sets the HTTP status code for the error
func (e *ManagedError) WithStatusCode(code int) *ManagedError {
	e.StatusCode = code
	return e
}

// IsType checks if the error is of a specific type
func IsType(err error, errType ErrorType) bool {
	var managedErr *ManagedError
	if errors.As(err, &managedErr) {
		return managedErr.Type == errType
	}
	return false
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	var managedErr *ManagedError
	if errors.As(err, &managedErr) {
		return managedErr.Retryable
	}
	return false
}

// GetContext retrieves context from an error
func GetContext(err error) map[string]string {
	var managedErr *ManagedError
	if errors.As(err, &managedErr) {
		return managedErr.Context
	}
	return nil
}

// Wrap wraps an existing error with additional context
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}
