// Package errmgt provides error management utilities for Go applications.
//
// This package offers structured error handling, error wrapping, and
// error categorization to help build robust applications with clear
// error reporting and debugging capabilities.
package errmgt

import (
	"fmt"
)

// ErrorType represents the category of an error
type ErrorType int

const (
	// ValidationError represents input validation errors
	ValidationError ErrorType = iota
	// NotFoundError represents resource not found errors
	NotFoundError
	// PermissionError represents authorization/permission errors
	PermissionError
	// InternalError represents internal system errors
	InternalError
	// ExternalError represents errors from external services
	ExternalError
)

// String returns the string representation of ErrorType
func (et ErrorType) String() string {
	switch et {
	case ValidationError:
		return "ValidationError"
	case NotFoundError:
		return "NotFoundError"
	case PermissionError:
		return "PermissionError"
	case InternalError:
		return "InternalError"
	case ExternalError:
		return "ExternalError"
	default:
		return "UnknownError"
	}
}

// ManagedError represents a structured error with type and context
type ManagedError struct {
	Type    ErrorType
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (me *ManagedError) Error() string {
	if me.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", me.Type, me.Message, me.Cause)
	}
	return fmt.Sprintf("[%s] %s", me.Type, me.Message)
}

// Unwrap returns the underlying cause error
func (me *ManagedError) Unwrap() error {
	return me.Cause
}

// New creates a new ManagedError with the specified type and message
func New(errorType ErrorType, message string) *ManagedError {
	return &ManagedError{
		Type:    errorType,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional context and type
func Wrap(err error, errorType ErrorType, message string) *ManagedError {
	return &ManagedError{
		Type:    errorType,
		Message: message,
		Cause:   err,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context information to the error
func (me *ManagedError) WithContext(key string, value interface{}) *ManagedError {
	me.Context[key] = value
	return me
}

// GetContext retrieves context information from the error
func (me *ManagedError) GetContext(key string) (interface{}, bool) {
	value, exists := me.Context[key]
	return value, exists
}

// IsType checks if the error is of a specific type
func (me *ManagedError) IsType(errorType ErrorType) bool {
	return me.Type == errorType
}
