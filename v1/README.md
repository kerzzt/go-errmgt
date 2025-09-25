# go-errmgt

[![CI](https://github.com/kerzzt/go-errmgt/actions/workflows/ci.yml/badge.svg)](https://github.com/kerzzt/go-errmgt/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kerzzt/go-errmgt.svg)](https://pkg.go.dev/github.com/kerzzt/go-errmgt)
[![Go Report Card](https://goreportcard.com/badge/github.com/kerzzt/go-errmgt)](https://goreportcard.com/report/github.com/kerzzt/go-errmgt)
[![Coverage](https://codecov.io/gh/kerzzt/go-errmgt/branch/main/graph/badge.svg)](https://codecov.io/gh/kerzzt/go-errmgt)

A comprehensive Golang library for structured error management and handling.

## Features

- **Structured Errors**: Create errors with type, code, message, and additional context
- **Error Categorization**: Built-in error types (Validation, Business, System, External)
- **Rich Context**: Add contextual information to errors for better debugging
- **Error Wrapping**: Compatible with Go's standard error wrapping patterns
- **Retry Logic Support**: Mark errors as retryable for automated retry mechanisms
- **HTTP Status Codes**: Associate HTTP status codes with errors for web applications
- **JSON Serializable**: Structured errors can be easily serialized to JSON

## Installation

```bash
go get github.com/kerzzt/go-errmgt
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/kerzzt/go-errmgt"
)

func main() {
    // Create a simple validation error
    err := errmgt.NewError(errmgt.ValidationError, "invalid_email", "Email format is invalid").
        WithDetails("Email must contain @ symbol").
        WithContext("field", "email").
        WithStatusCode(400)

    fmt.Printf("Error: %s\n", err.Error())
    // Output: Error: [validation:invalid_email] Email format is invalid: Email must contain @ symbol
}
```

## Error Types

The library provides four built-in error types:

- `ValidationError`: Input validation errors
- `BusinessError`: Business logic errors  
- `SystemError`: System-level errors (database, network, etc.)
- `ExternalError`: Errors from external services

## API Reference

### Creating Errors

```go
// Simple error
err := errmgt.NewError(errmgt.ValidationError, "invalid_input", "Input validation failed")

// Error with cause
dbErr := fmt.Errorf("connection timeout")
err := errmgt.NewErrorWithCause(errmgt.SystemError, "db_error", "Database operation failed", dbErr)
```

### Adding Context

```go
err := errmgt.NewError(errmgt.BusinessError, "insufficient_funds", "Not enough balance").
    WithDetails("Account balance: $10, Required: $50").
    WithContext("user_id", "12345").
    WithContext("account_id", "67890").
    WithRetryable(false).
    WithStatusCode(402)
```

### Error Inspection

```go
// Check error type
if errmgt.IsType(err, errmgt.ValidationError) {
    // Handle validation error
}

// Check if retryable
if errmgt.IsRetryable(err) {
    // Retry the operation
}

// Get context
context := errmgt.GetContext(err)
userID := context["user_id"]
```

### Error Wrapping

```go
// Simple wrapping
wrappedErr := errmgt.Wrap(originalErr, "additional context")

// Formatted wrapping
wrappedErr := errmgt.Wrapf(originalErr, "failed to process user %s", userID)
```

## Usage Examples

### Web Application Error Handling

```go
func registerUser(email, password string) error {
    if !isValidEmail(email) {
        return errmgt.NewError(errmgt.ValidationError, "invalid_email", "Invalid email format").
            WithDetails("Must be a valid email address").
            WithContext("email", email).
            WithStatusCode(400)
    }
    
    if err := saveUser(email, password); err != nil {
        return errmgt.NewErrorWithCause(errmgt.SystemError, "db_save_failed", "Failed to save user", err).
            WithRetryable(true).
            WithContext("operation", "user_registration").
            WithStatusCode(500)
    }
    
    return nil
}

func handleError(w http.ResponseWriter, err error) {
    var managedErr *errmgt.ManagedError
    if errors.As(err, &managedErr) {
        statusCode := managedErr.StatusCode
        if statusCode == 0 {
            statusCode = 500 // default
        }
        
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(managedErr)
    } else {
        w.WriteHeader(500)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
    }
}
```

### Retry Logic

```go
func processWithRetry(fn func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        if err := fn(); err != nil {
            lastErr = err
            if errmgt.IsRetryable(err) {
                time.Sleep(time.Duration(i+1) * time.Second)
                continue
            }
            return err // Non-retryable error
        }
        return nil // Success
    }
    
    return errmgt.Wrapf(lastErr, "operation failed after %d retries", maxRetries)
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.