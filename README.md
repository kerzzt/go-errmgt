# go-errmgt

[![CI](https://github.com/kerzzt/go-errmgt/actions/workflows/ci.yml/badge.svg)](https://github.com/kerzzt/go-errmgt/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kerzzt/go-errmgt.svg)](https://pkg.go.dev/github.com/kerzzt/go-errmgt)
[![Go Report Card](https://goreportcard.com/badge/github.com/kerzzt/go-errmgt)](https://goreportcard.com/report/github.com/kerzzt/go-errmgt)

A comprehensive Golang library for structured error management and handling.

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
    err := errmgt.NewError(errmgt.ValidationError, "invalid_email", "Email format is invalid").
        WithDetails("Email must contain @ symbol").
        WithContext("field", "email").
        WithStatusCode(400)

    fmt.Printf("Error: %s\n", err.Error())
}
```

## Documentation

For complete documentation, examples, and API reference, see the [v1 directory](./v1/README.md).

## Features

- **Structured Errors**: Create errors with type, code, message, and additional context
- **Error Categorization**: Built-in error types (Validation, Business, System, External)
- **Rich Context**: Add contextual information to errors for better debugging
- **Error Wrapping**: Compatible with Go's standard error wrapping patterns
- **Retry Logic Support**: Mark errors as retryable for automated retry mechanisms
- **HTTP Status Codes**: Associate HTTP status codes with errors for web applications
- **JSON Serializable**: Structured errors can be easily serialized to JSON

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.