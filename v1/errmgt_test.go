package errmgt

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(ValidationError, "invalid_input", "Input validation failed")
	
	if err.Type != ValidationError {
		t.Errorf("Expected type %s, got %s", ValidationError, err.Type)
	}
	
	if err.Code != "invalid_input" {
		t.Errorf("Expected code 'invalid_input', got '%s'", err.Code)
	}
	
	if err.Message != "Input validation failed" {
		t.Errorf("Expected message 'Input validation failed', got '%s'", err.Message)
	}
	
	if err.Context == nil {
		t.Error("Expected context to be initialized")
	}
}

func TestNewErrorWithCause(t *testing.T) {
	cause := errors.New("original error")
	err := NewErrorWithCause(SystemError, "db_connection", "Database connection failed", cause)
	
	if err.Type != SystemError {
		t.Errorf("Expected type %s, got %s", SystemError, err.Type)
	}
	
	if err.Cause != cause {
		t.Error("Expected cause to be set")
	}
	
	if !errors.Is(err, cause) {
		t.Error("Expected error to be identified as the cause")
	}
}

func TestManagedErrorError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ManagedError
		expected string
	}{
		{
			name: "without details",
			err:  NewError(ValidationError, "invalid_email", "Invalid email format"),
			expected: "[validation:invalid_email] Invalid email format",
		},
		{
			name: "with details",
			err:  NewError(ValidationError, "invalid_email", "Invalid email format").WithDetails("Email must contain @ symbol"),
			expected: "[validation:invalid_email] Invalid email format: Email must contain @ symbol",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestManagedErrorWithMethods(t *testing.T) {
	err := NewError(BusinessError, "insufficient_funds", "Insufficient account balance")
	
	// Test WithDetails
	err = err.WithDetails("Account balance: $10, Required: $50")
	if err.Details != "Account balance: $10, Required: $50" {
		t.Errorf("Expected details to be set, got '%s'", err.Details)
	}
	
	// Test WithContext
	err = err.WithContext("user_id", "12345").WithContext("account_id", "67890")
	if err.Context["user_id"] != "12345" {
		t.Error("Expected user_id context to be set")
	}
	if err.Context["account_id"] != "67890" {
		t.Error("Expected account_id context to be set")
	}
	
	// Test WithRetryable
	err = err.WithRetryable(true)
	if !err.Retryable {
		t.Error("Expected error to be retryable")
	}
	
	// Test WithStatusCode
	err = err.WithStatusCode(402)
	if err.StatusCode != 402 {
		t.Errorf("Expected status code 402, got %d", err.StatusCode)
	}
}

func TestIsType(t *testing.T) {
	validationErr := NewError(ValidationError, "invalid_input", "Invalid input")
	businessErr := NewError(BusinessError, "business_rule", "Business rule violation")
	regularErr := errors.New("regular error")
	
	if !IsType(validationErr, ValidationError) {
		t.Error("Expected validation error to be identified as ValidationError")
	}
	
	if IsType(validationErr, BusinessError) {
		t.Error("Expected validation error not to be identified as BusinessError")
	}
	
	if IsType(regularErr, ValidationError) {
		t.Error("Expected regular error not to be identified as ValidationError")
	}
	
	if !IsType(businessErr, BusinessError) {
		t.Error("Expected business error to be identified as BusinessError")
	}
}

func TestIsRetryable(t *testing.T) {
	retryableErr := NewError(ExternalError, "api_timeout", "API timeout").WithRetryable(true)
	nonRetryableErr := NewError(ValidationError, "invalid_input", "Invalid input").WithRetryable(false)
	regularErr := errors.New("regular error")
	
	if !IsRetryable(retryableErr) {
		t.Error("Expected retryable error to be identified as retryable")
	}
	
	if IsRetryable(nonRetryableErr) {
		t.Error("Expected non-retryable error not to be identified as retryable")
	}
	
	if IsRetryable(regularErr) {
		t.Error("Expected regular error not to be identified as retryable")
	}
}

func TestGetContext(t *testing.T) {
	err := NewError(SystemError, "db_error", "Database error").
		WithContext("table", "users").
		WithContext("operation", "select")
	
	context := GetContext(err)
	if context == nil {
		t.Fatal("Expected context to be returned")
	}
	
	if context["table"] != "users" {
		t.Error("Expected table context to be 'users'")
	}
	
	if context["operation"] != "select" {
		t.Error("Expected operation context to be 'select'")
	}
	
	regularErr := errors.New("regular error")
	context = GetContext(regularErr)
	if context != nil {
		t.Error("Expected no context for regular error")
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, "additional context")
	
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Expected wrapped error to be identified as original error")
	}
	
	expectedMsg := "additional context: original error"
	if wrappedErr.Error() != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, wrappedErr.Error())
	}
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("connection failed")
	wrappedErr := Wrapf(originalErr, "failed to connect to %s:%d", "localhost", 5432)
	
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Expected wrapped error to be identified as original error")
	}
	
	expectedMsg := "failed to connect to localhost:5432: connection failed"
	if wrappedErr.Error() != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, wrappedErr.Error())
	}
}

func TestManagedErrorIs(t *testing.T) {
	// Test with same type and code
	err1 := NewError(ValidationError, "invalid_email", "Invalid email")
	err2 := NewError(ValidationError, "invalid_email", "Different message")
	
	if !errors.Is(err1, err2) {
		t.Error("Expected errors with same type and code to be equal")
	}
	
	// Test with different type
	err3 := NewError(BusinessError, "invalid_email", "Invalid email")
	if errors.Is(err1, err3) {
		t.Error("Expected errors with different types not to be equal")
	}
	
	// Test with different code
	err4 := NewError(ValidationError, "invalid_phone", "Invalid phone")
	if errors.Is(err1, err4) {
		t.Error("Expected errors with different codes not to be equal")
	}
	
	// Test with underlying cause
	cause := errors.New("underlying error")
	err5 := NewErrorWithCause(SystemError, "db_error", "Database error", cause)
	
	if !errors.Is(err5, cause) {
		t.Error("Expected error to be identified as its cause")
	}
}