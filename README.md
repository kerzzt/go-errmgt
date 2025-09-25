# go-errmgt

Error management Golang library that provides structured error handling, error wrapping, and error categorization to help build robust applications with clear error reporting and debugging capabilities.

## Features

- **Structured Error Types**: Categorize errors into predefined types (Validation, NotFound, Permission, Internal, External)
- **Error Wrapping**: Wrap existing errors with additional context while preserving the original error chain
- **Context Management**: Add and retrieve contextual information to/from errors
- **Go Standard Library Compatible**: Works seamlessly with `errors.Is` and `errors.Unwrap`
- **Type-Safe Error Checking**: Check error types with compile-time safety

## Installation

```bash
go get github.com/kerzzt/go-errmgt
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/kerzzt/go-errmgt"
)

func main() {
    // Create a new structured error
    err := errmgt.New(errmgt.ValidationError, "email is required")
    fmt.Println(err.Error()) // [ValidationError] email is required
    
    // Add context to errors
    err.WithContext("field", "email").WithContext("userId", 123)
    
    // Check error type
    if err.IsType(errmgt.ValidationError) {
        log.Println("Handle validation error")
    }
}
```

## Error Types

The library provides five predefined error types:

- `ValidationError` - Input validation errors
- `NotFoundError` - Resource not found errors  
- `PermissionError` - Authorization/permission errors
- `InternalError` - Internal system errors
- `ExternalError` - Errors from external services

## API Reference

### Creating Errors

#### `New(errorType ErrorType, message string) *ManagedError`
Creates a new structured error with the specified type and message.

```go
err := errmgt.New(errmgt.ValidationError, "invalid email format")
```

#### `Wrap(err error, errorType ErrorType, message string) *ManagedError`
Wraps an existing error with additional context and type information.

```go
originalErr := errors.New("connection timeout")
wrappedErr := errmgt.Wrap(originalErr, errmgt.ExternalError, "failed to connect to API")
```

### Working with Context

#### `WithContext(key string, value interface{}) *ManagedError`
Adds contextual information to the error.

```go
err := errmgt.New(errmgt.ValidationError, "invalid input").
    WithContext("field", "email").
    WithContext("userId", 123)
```

#### `GetContext(key string) (interface{}, bool)`
Retrieves contextual information from the error.

```go
if field, exists := err.GetContext("field"); exists {
    fmt.Printf("Field: %s\n", field)
}
```

### Error Checking

#### `IsType(errorType ErrorType) bool`
Checks if the error is of a specific type.

```go
if err.IsType(errmgt.ValidationError) {
    // Handle validation error
}
```

#### `Error() string`
Returns the string representation of the error.

#### `Unwrap() error`
Returns the underlying cause error (compatible with `errors.Unwrap`).

## Examples

### Basic Usage

```go
// Create and handle different error types
func processUser(email string) error {
    if email == "" {
        return errmgt.New(errmgt.ValidationError, "email is required")
    }
    
    // ... other validation logic
    
    return nil
}

func handleError(err error) {
    if managedErr, ok := err.(*errmgt.ManagedError); ok {
        switch {
        case managedErr.IsType(errmgt.ValidationError):
            log.Printf("Validation failed: %s", managedErr.Error())
        case managedErr.IsType(errmgt.NotFoundError):
            log.Printf("Resource not found: %s", managedErr.Error())
        default:
            log.Printf("Other error: %s", managedErr.Error())
        }
    }
}
```

### Error Wrapping

```go
func connectToDatabase() error {
    conn, err := sql.Open("postgres", "connection-string")
    if err != nil {
        return errmgt.Wrap(err, errmgt.InternalError, "failed to connect to database").
            WithContext("database", "postgres").
            WithContext("operation", "connect")
    }
    return nil
}
```

### Integration with Standard Library

```go
func example() {
    originalErr := errors.New("network error")
    wrappedErr := errmgt.Wrap(originalErr, errmgt.ExternalError, "API call failed")
    
    // Standard library compatibility
    if errors.Is(wrappedErr, originalErr) {
        fmt.Println("Original error found in chain")
    }
    
    if unwrapped := errors.Unwrap(wrappedErr); unwrapped != nil {
        fmt.Printf("Unwrapped: %s\n", unwrapped.Error())
    }
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
