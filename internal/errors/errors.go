package errors

import (
	"fmt"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// Validation errors
	ErrorTypeValidation ErrorType = "validation"

	// Storage errors
	ErrorTypeStorage ErrorType = "storage"

	// Business logic errors
	ErrorTypeBusiness ErrorType = "business"

	// CLI/User input errors
	ErrorTypeUserInput ErrorType = "user_input"

	// System errors
	ErrorTypeSystem ErrorType = "system"
)

// ComptesError represents a custom error with context
type ComptesError struct {
	Type    ErrorType
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *ComptesError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Type, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *ComptesError) Unwrap() error {
	return e.Cause
}

// New creates a new ComptesError
func New(errorType ErrorType, code, message string) *ComptesError {
	return &ComptesError{
		Type:    errorType,
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with context
func Wrap(errorType ErrorType, code, message string, cause error) *ComptesError {
	return &ComptesError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Error codes
const (
	// Validation error codes
	CodeAccountNotFound     = "account_not_found"
	CodeCategoryNotFound    = "category_not_found"
	CodeTagNotFound         = "tag_not_found"
	CodeTransactionNotFound = "transaction_not_found"
	CodeInvalidAmount       = "invalid_amount"
	CodeInvalidDate         = "invalid_date"
	CodeMissingField        = "missing_field"

	// Storage error codes
	CodeStorageReadFailed  = "storage_read_failed"
	CodeStorageWriteFailed = "storage_write_failed"
	CodeStorageInitFailed  = "storage_init_failed"

	// Business logic error codes
	CodeTransactionAlreadyDeleted = "transaction_already_deleted"
	CodeInvalidOperation          = "invalid_operation"
	CodeAmbiguousID               = "ambiguous_id"
	CodeParentNotFound            = "parent_not_found"

	// User input error codes
	CodeMissingArguments = "missing_arguments"
	CodeMissingMessage   = "missing_message"
	CodeInvalidJSON      = "invalid_json"
	CodeInvalidCommand   = "invalid_command"

	// System error codes
	CodeConfigLoadFailed = "config_load_failed"
	CodeConfigSaveFailed = "config_save_failed"
)

// Predefined error constructors for common cases

// Validation errors
func AccountNotFound(accountID string) *ComptesError {
	return New(ErrorTypeValidation, CodeAccountNotFound, fmt.Sprintf("Account not found: %s", accountID))
}

func CategoryNotFound(categoryCode string) *ComptesError {
	return New(ErrorTypeValidation, CodeCategoryNotFound, fmt.Sprintf("Category not found: %s", categoryCode))
}

func TagNotFound(tagCode string) *ComptesError {
	return New(ErrorTypeValidation, CodeTagNotFound, fmt.Sprintf("Tag not found: %s", tagCode))
}

func TransactionNotFound(transactionID string) *ComptesError {
	return New(ErrorTypeValidation, CodeTransactionNotFound, fmt.Sprintf("Transaction not found: %s", transactionID))
}

func AmbiguousID(partialID string) *ComptesError {
	return New(ErrorTypeValidation, CodeAmbiguousID, fmt.Sprintf("Multiple transactions found with ID starting with: %s (be more specific)", partialID))
}

// Storage errors
func StorageReadFailed(resource string, cause error) *ComptesError {
	return Wrap(ErrorTypeStorage, CodeStorageReadFailed, fmt.Sprintf("Failed to read %s", resource), cause)
}

func StorageWriteFailed(resource string, cause error) *ComptesError {
	return Wrap(ErrorTypeStorage, CodeStorageWriteFailed, fmt.Sprintf("Failed to write %s", resource), cause)
}

// Business logic errors
func TransactionAlreadyDeleted(transactionID string) *ComptesError {
	return New(ErrorTypeBusiness, CodeTransactionAlreadyDeleted, fmt.Sprintf("Transaction %s is already deleted", transactionID))
}

func InvalidOperation(transactionID string) *ComptesError {
	return New(ErrorTypeBusiness, CodeInvalidOperation, fmt.Sprintf("Cannot determine operation type for transaction %s", transactionID))
}

func ParentNotFound(parentID string) *ComptesError {
	return New(ErrorTypeBusiness, CodeParentNotFound, fmt.Sprintf("Parent transaction %s not found", parentID))
}

// User input errors
func MissingArguments(command string) *ComptesError {
	return New(ErrorTypeUserInput, CodeMissingArguments, fmt.Sprintf("Missing arguments for command: %s", command))
}

func MissingMessage(operation string) *ComptesError {
	return New(ErrorTypeUserInput, CodeMissingMessage, fmt.Sprintf("Message is mandatory for %s operations", operation))
}

func InvalidJSON(cause error) *ComptesError {
	return Wrap(ErrorTypeUserInput, CodeInvalidJSON, "Failed to parse JSON", cause)
}

func InvalidCommand(command string) *ComptesError {
	return New(ErrorTypeUserInput, CodeInvalidCommand, fmt.Sprintf("Unknown command: %s", command))
}

// System errors
func ConfigLoadFailed(cause error) *ComptesError {
	return Wrap(ErrorTypeSystem, CodeConfigLoadFailed, "Failed to load configuration", cause)
}

func ConfigSaveFailed(cause error) *ComptesError {
	return Wrap(ErrorTypeSystem, CodeConfigSaveFailed, "Failed to save configuration", cause)
}
