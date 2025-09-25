package errmgt_test

import (
	"errors"
	"fmt"

	"github.com/kerzzt/go-errmgt"
)

func ExampleNew() {
	// Create a new validation error
	err := errmgt.New(errmgt.ValidationError, "email is required")
	fmt.Println(err.Error())
	// Output: [ValidationError] email is required
}

func ExampleWrap() {
	// Wrap an existing error with additional context
	originalErr := errors.New("connection timeout")
	wrappedErr := errmgt.Wrap(originalErr, errmgt.ExternalError, "failed to connect to API")
	fmt.Println(wrappedErr.Error())
	// Output: [ExternalError] failed to connect to API: connection timeout
}

func ExampleManagedError_WithContext() {
	// Create an error with additional context information
	err := errmgt.New(errmgt.ValidationError, "invalid input").
		WithContext("field", "email").
		WithContext("userId", 123)

	fmt.Println(err.Error())

	// Retrieve context
	if field, exists := err.GetContext("field"); exists {
		fmt.Printf("Field: %s\n", field)
	}
	// Output:
	// [ValidationError] invalid input
	// Field: email
}

func ExampleManagedError_IsType() {
	err := errmgt.New(errmgt.NotFoundError, "user not found")

	if err.IsType(errmgt.NotFoundError) {
		fmt.Println("This is a not found error")
	}
	// Output: This is a not found error
}

// Example of using the error management library in a web service context
func Example_webServiceUsage() {
	// Simulate a service function that might encounter various error types
	validateUser := func(email string) error {
		if email == "" {
			return errmgt.New(errmgt.ValidationError, "email is required").
				WithContext("field", "email")
		}
		if email == "blocked@example.com" {
			return errmgt.New(errmgt.PermissionError, "user is blocked").
				WithContext("email", email)
		}
		return nil
	}

	// Handle different error types
	handleError := func(err error) {
		if managedErr, ok := err.(*errmgt.ManagedError); ok {
			switch {
			case managedErr.IsType(errmgt.ValidationError):
				fmt.Printf("Validation error: %s\n", managedErr.Error())
			case managedErr.IsType(errmgt.PermissionError):
				fmt.Printf("Permission error: %s\n", managedErr.Error())
			default:
				fmt.Printf("Other error: %s\n", managedErr.Error())
			}
		}
	}

	// Test cases
	if err := validateUser(""); err != nil {
		handleError(err)
	}

	if err := validateUser("blocked@example.com"); err != nil {
		handleError(err)
	}

	// Output:
	// Validation error: [ValidationError] email is required
	// Permission error: [PermissionError] user is blocked
}
