package database

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Severity represents the PostgreSQL error severity level
type Severity string

// PostgreSQL severity levels
const (
	SeverityError   Severity = "ERROR"
	SeverityFatal   Severity = "FATAL"
	SeverityPanic   Severity = "PANIC"
	SeverityWarning Severity = "WARNING"
	SeverityNotice  Severity = "NOTICE"
	SeverityDebug   Severity = "DEBUG"
	SeverityInfo    Severity = "INFO"
	SeverityLog     Severity = "LOG"
)

// ErrorClass represents the classification of an error
type ErrorClass int

// Error classifications
const (
	// ClassUnknown represents an unknown error
	ClassUnknown ErrorClass = iota
	// ClassConnection represents connection-related errors
	ClassConnection
	// ClassConstraint represents constraint violation errors
	ClassConstraint
	// ClassDataException represents data exceptions
	ClassDataException
	// ClassTransactionFailure represents transaction failures
	ClassTransactionFailure
	// ClassSyntaxError represents syntax errors
	ClassSyntaxError
	// ClassResourceError represents resource errors (out of memory, etc.)
	ClassResourceError
	// ClassSystemError represents internal system errors
	ClassSystemError
	// ClassConfiguration represents configuration errors
	ClassConfiguration
)

// ErrorInfo contains detailed information about a PostgreSQL error
type ErrorInfo struct {
	Err *pgconn.PgError
	// The PostgreSQL error code
	Code string
	// The severity level
	Severity Severity
	// The error class
	Class ErrorClass
	// Whether the error can potentially be retried
	Retryable bool
	// Whether the application can handle this error
	Recoverable bool
	// Human-readable message
	Message string
	// Additional details
	Detail string
	// Hint for solving the problem
	Hint string
}

// Error implements the error interface
func (e ErrorInfo) Error() string {
	return fmt.Sprintf("%s: %s (Code: %s, Severity: %s, Class: %d, Retryable: %t, Recoverable: %t)",
		e.Severity, e.Message, e.Code, e.Severity, e.Class, e.Retryable, e.Recoverable)
}

// Unwrap returns the original error
func (e ErrorInfo) Unwrap() error {
	return e.Err
}

// ClassifyError examines a PostgreSQL error and returns detailed information
func ClassifyError(err error) *ErrorInfo {
	if err == nil {
		return nil
	}

	// Handle pgx.ErrNoRows specifically
	if errors.Is(err, pgx.ErrNoRows) {
		return &ErrorInfo{
			Err:         nil,
			Code:        "P0002", // No data (client-defined code)
			Severity:    SeverityInfo,
			Class:       ClassDataException,
			Retryable:   false,
			Recoverable: true,
			Message:     "No rows found",
		}
	}

	// Extract the PostgreSQL error if available
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		info := &ErrorInfo{
			Err:      pgErr,
			Code:     pgErr.Code,
			Severity: Severity(pgErr.Severity),
			Message:  pgErr.Message,
			Detail:   pgErr.Detail,
			Hint:     pgErr.Hint,
		}

		// Classify based on error code class (first two characters)
		codeClass := pgErr.Code[:2]
		switch codeClass {
		case "08": // Connection Exception
			info.Class = ClassConnection
			info.Retryable = true
			info.Recoverable = false
		case "23": // Integrity Constraint Violation
			info.Class = ClassConstraint
			info.Retryable = false
			info.Recoverable = true
		case "22": // Data Exception
			info.Class = ClassDataException
			info.Retryable = false
			info.Recoverable = true
		case "40": // Transaction Rollback
			info.Class = ClassTransactionFailure
			info.Retryable = true
			info.Recoverable = true
		case "42": // Syntax Error or Access Rule Violation
			info.Class = ClassSyntaxError
			info.Retryable = false
			info.Recoverable = false
		case "53": // Insufficient Resources
			info.Class = ClassResourceError
			info.Retryable = true
			info.Recoverable = false
		case "XX": // Internal Error
			info.Class = ClassSystemError
			info.Retryable = false
			info.Recoverable = false
		case "F0": // Configuration File Error
			info.Class = ClassConfiguration
			info.Retryable = false
			info.Recoverable = false
		default:
			info.Class = ClassUnknown
			info.Retryable = false
			info.Recoverable = false
		}

		// Adjust based on severity
		switch info.Severity {
		case SeverityError:
			// Most ERROR level issues can be handled by the application
			// But might need to override based on specific codes
			if info.Class == ClassUnknown {
				info.Recoverable = true
			}
		case SeverityFatal, SeverityPanic:
			// FATAL and PANIC errors indicate serious problems
			info.Retryable = false
			info.Recoverable = false
		case SeverityWarning, SeverityNotice, SeverityDebug, SeverityInfo, SeverityLog:
			// These are informational and don't indicate actual errors
			info.Retryable = true
			info.Recoverable = true
		}

		// Special case handling for specific error codes
		switch pgErr.Code {
		case "40001", "40P01": // Serialization failure, deadlock detected
			info.Retryable = true
			info.Recoverable = true
		case "08006": // Connection failure
			info.Retryable = true
			info.Recoverable = false
		case "57P01", "57P02", "57P03": // Admin shutdown, crash shutdown, cannot connect now
			info.Retryable = true
			info.Recoverable = false
		case "53300": // Too many connections
			info.Retryable = true
			info.Recoverable = false
		}

		return info
	}

	// Non-PostgreSQL error
	return nil
}

// IsRetryable checks if an error can be retried
func (e *ErrorInfo) IsRetryable() bool {
	return e != nil && e.Retryable
}

// IsRecoverable checks if an error can be handled by the application
func (e *ErrorInfo) IsRecoverable() bool {
	return e != nil && e.Recoverable
}

// IsFatalError checks if an error is fatal
func (e *ErrorInfo) IsFatalError() bool {
	return e != nil && (e.Severity == SeverityFatal || e.Severity == SeverityPanic)
}

// IsConstraintViolation checks if an error is a constraint violation
func (e *ErrorInfo) IsConstraintViolation() bool {
	return e != nil && e.Class == ClassConstraint
}

// GetErrorClass returns the error class
func (e *ErrorInfo) GetErrorClass() ErrorClass {
	if e == nil {
		return ClassUnknown
	}
	return e.Class
}

// GetErrorSeverity returns the error severity
func (e *ErrorInfo) GetErrorSeverity() Severity {
	if e == nil {
		return SeverityError
	}
	return e.Severity
}
