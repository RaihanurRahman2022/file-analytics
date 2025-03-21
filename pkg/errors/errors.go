package errors

import (
	"errors"
	"fmt"
	"time"
)

// ErrorType represents different types of errors
// Demonstrates enum pattern with iota
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeIO
	ErrorTypeFormat
	ErrorTypeTimeout
	ErrorTypeValidation
)

// String implements Stringer interface
func (et ErrorType) String() string {
	// Demonstrates switch statement
	switch et {
	case ErrorTypeIO:
		return "IO Error"
	case ErrorTypeFormat:
		return "Format Error"
	case ErrorTypeTimeout:
		return "Timeout Error"
	case ErrorTypeValidation:
		return "Validation Error"
	default:
		return "Unknown Error"
	}
}

// ProcessError represents an error that occurred during file processing
// Demonstrates custom error type
type ProcessError struct {
	Type    ErrorType
	File    string
	Message string
	Cause   error
	Time    time.Time
}

// Error implements the error interface
// Demonstrates error interface implementation
func (e *ProcessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s in file '%s': %v", e.Type, e.Message, e.File, e.Cause)
	}
	return fmt.Sprintf("%s: %s in file '%s'", e.Type, e.Message, e.File)
}

// Unwrap implements error unwrapping
// Demonstrates error wrapping
func (e *ProcessError) Unwrap() error {
	return e.Cause
}

// NewProcessError creates a new ProcessError
// Demonstrates variadic error constructor
func NewProcessError(errType ErrorType, file string, message string, causes ...error) *ProcessError {
	var cause error
	if len(causes) > 0 {
		cause = causes[0]
	}

	return &ProcessError{
		Type:    errType,
		File:    file,
		Message: message,
		Cause:   cause,
		Time:    time.Now(),
	}
}

// ErrorCollection represents a collection of errors
// Demonstrates slice usage with errors
type ErrorCollection struct {
	errors []error
}

// NewErrorCollection creates a new error collection
func NewErrorCollection() *ErrorCollection {
	return &ErrorCollection{
		errors: make([]error, 0),
	}
}

// Add adds an error to the collection
// Demonstrates pointer receiver method
func (ec *ErrorCollection) Add(err error) {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
}

// HasErrors checks if the collection contains any errors
func (ec *ErrorCollection) HasErrors() bool {
	return len(ec.errors) > 0
}

// Errors returns all errors in the collection
// Demonstrates slice return
func (ec *ErrorCollection) Errors() []error {
	return ec.errors
}

// Error implements the error interface
// Demonstrates string building
func (ec *ErrorCollection) Error() string {
	if !ec.HasErrors() {
		return "no errors"
	}

	result := fmt.Sprintf("%d error(s) occurred:\n", len(ec.errors))
	for i, err := range ec.errors {
		result += fmt.Sprintf("%d. %v\n", i+1, err)
	}
	return result
}

// IsErrorType checks if an error is of a specific type
// Demonstrates type assertion and error handling
func IsErrorType(err error, errType ErrorType) bool {
	var processErr *ProcessError
	if ok := errors.As(err, &processErr); ok {
		return processErr.Type == errType
	}
	return false
}

// As attempts to convert an error to a specific type
// Demonstrates generic error handling
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Wrap wraps an error with additional context
// Demonstrates error wrapping utility
func Wrap(err error, errType ErrorType, file string, message string) error {
	if err == nil {
		return nil
	}
	return NewProcessError(errType, file, message, err)
}
