package errmgt

import (
	"errors"
	"testing"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{ValidationError, "ValidationError"},
		{NotFoundError, "NotFoundError"},
		{PermissionError, "PermissionError"},
		{InternalError, "InternalError"},
		{ExternalError, "ExternalError"},
		{ErrorType(999), "UnknownError"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.errorType.String(); got != test.expected {
				t.Errorf("ErrorType.String() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	err := New(ValidationError, "test validation error")

	if err.Type != ValidationError {
		t.Errorf("New() Type = %v, want %v", err.Type, ValidationError)
	}

	if err.Message != "test validation error" {
		t.Errorf("New() Message = %v, want %v", err.Message, "test validation error")
	}

	if err.Cause != nil {
		t.Errorf("New() Cause = %v, want nil", err.Cause)
	}

	if err.Context == nil {
		t.Error("New() Context should be initialized")
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, InternalError, "wrapped error message")

	if wrappedErr.Type != InternalError {
		t.Errorf("Wrap() Type = %v, want %v", wrappedErr.Type, InternalError)
	}

	if wrappedErr.Message != "wrapped error message" {
		t.Errorf("Wrap() Message = %v, want %v", wrappedErr.Message, "wrapped error message")
	}

	if wrappedErr.Cause != originalErr {
		t.Errorf("Wrap() Cause = %v, want %v", wrappedErr.Cause, originalErr)
	}
}

func TestManagedError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ManagedError
		expected string
	}{
		{
			name:     "error without cause",
			err:      New(ValidationError, "validation failed"),
			expected: "[ValidationError] validation failed",
		},
		{
			name:     "error with cause",
			err:      Wrap(errors.New("original"), InternalError, "internal error"),
			expected: "[InternalError] internal error: original",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.err.Error(); got != test.expected {
				t.Errorf("ManagedError.Error() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestManagedError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, InternalError, "wrapped")

	if unwrapped := wrappedErr.Unwrap(); unwrapped != originalErr {
		t.Errorf("ManagedError.Unwrap() = %v, want %v", unwrapped, originalErr)
	}

	// Test unwrapping nil cause
	newErr := New(ValidationError, "test")
	if unwrapped := newErr.Unwrap(); unwrapped != nil {
		t.Errorf("ManagedError.Unwrap() = %v, want nil for error without cause", unwrapped)
	}
}

func TestManagedError_WithContext(t *testing.T) {
	err := New(ValidationError, "test error")
	err.WithContext("userId", 123)
	err.WithContext("field", "email")

	if value, exists := err.GetContext("userId"); !exists || value != 123 {
		t.Errorf("WithContext/GetContext userId = %v, %v, want 123, true", value, exists)
	}

	if value, exists := err.GetContext("field"); !exists || value != "email" {
		t.Errorf("WithContext/GetContext field = %v, %v, want 'email', true", value, exists)
	}

	if _, exists := err.GetContext("nonexistent"); exists {
		t.Error("GetContext should return false for non-existent key")
	}
}

func TestManagedError_IsType(t *testing.T) {
	err := New(ValidationError, "test error")

	if !err.IsType(ValidationError) {
		t.Error("IsType() should return true for matching error type")
	}

	if err.IsType(NotFoundError) {
		t.Error("IsType() should return false for non-matching error type")
	}
}

func TestErrorsIs(t *testing.T) {
	// Test compatibility with errors.Is
	originalErr := errors.New("original")
	wrappedErr := Wrap(originalErr, InternalError, "wrapped")

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is should work with wrapped ManagedError")
	}
}
