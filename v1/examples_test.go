package errmgt_test

import (
	"fmt"

	"github.com/kerzzt/go-errmgt"
)

func ExampleNewError() {
	err := errmgt.NewError(errmgt.ValidationError, "invalid_email", "Email format is invalid").
		WithDetails("Email must contain @ symbol").
		WithContext("field", "email").
		WithStatusCode(400)

	fmt.Printf("Error: %s\n", err.Error())
	fmt.Printf("Type: %s\n", err.Type)
	fmt.Printf("Code: %s\n", err.Code)
	fmt.Printf("Status Code: %d\n", err.StatusCode)
	
	// Output:
	// Error: [validation:invalid_email] Email format is invalid: Email must contain @ symbol
	// Type: validation
	// Code: invalid_email
	// Status Code: 400
}

func ExampleIsType() {
	err := errmgt.NewError(errmgt.BusinessError, "insufficient_funds", "Not enough money")
	
	if errmgt.IsType(err, errmgt.BusinessError) {
		fmt.Println("This is a business logic error")
	}
	
	// Output:
	// This is a business logic error
}

func ExampleNewErrorWithCause() {
	// Simulate a database connection error
	dbErr := fmt.Errorf("connection timeout")
	
	err := errmgt.NewErrorWithCause(
		errmgt.SystemError,
		"db_connection_failed",
		"Failed to connect to database",
		dbErr,
	).WithRetryable(true).WithContext("database", "users_db")

	fmt.Printf("Error: %s\n", err.Error())
	fmt.Printf("Retryable: %t\n", errmgt.IsRetryable(err))
	
	// Output:
	// Error: [system:db_connection_failed] Failed to connect to database
	// Retryable: true
}

func ExampleWrap() {
	originalErr := fmt.Errorf("file not found")
	wrappedErr := errmgt.Wrap(originalErr, "failed to read configuration")
	
	fmt.Printf("Wrapped error: %s\n", wrappedErr.Error())
	
	// Output:
	// Wrapped error: failed to read configuration: file not found
}

// Example of how to use the library in a real application
func Example_usageInApplication() {
	// Simulate user registration
	err := registerUser("invalid-email", "password123")
	if err != nil {
		if errmgt.IsType(err, errmgt.ValidationError) {
			fmt.Printf("Validation failed: %s\n", err.Error())
			// Handle validation error (e.g., return 400 status)
		} else if errmgt.IsType(err, errmgt.SystemError) && errmgt.IsRetryable(err) {
			fmt.Printf("System error (retryable): %s\n", err.Error())
			// Handle retryable system error (e.g., retry operation)
		} else {
			fmt.Printf("Unexpected error: %s\n", err.Error())
		}
	}
	
	// Output:
	// Validation failed: [validation:invalid_email] Invalid email format: Must contain @ symbol
}

func registerUser(email, password string) error {
	// Validate email
	if !contains(email, "@") {
		return errmgt.NewError(errmgt.ValidationError, "invalid_email", "Invalid email format").
			WithDetails("Must contain @ symbol").
			WithContext("email", email).
			WithStatusCode(400)
	}
	
	// Validate password
	if len(password) < 8 {
		return errmgt.NewError(errmgt.ValidationError, "weak_password", "Password too weak").
			WithDetails("Password must be at least 8 characters").
			WithStatusCode(400)
	}
	
	// Simulate database save
	if err := saveToDatabase(email, password); err != nil {
		return errmgt.NewErrorWithCause(errmgt.SystemError, "db_save_failed", "Failed to save user", err).
			WithRetryable(true).
			WithContext("operation", "user_registration")
	}
	
	return nil
}

func saveToDatabase(email, password string) error {
	// Simulate database operation
	return nil
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}