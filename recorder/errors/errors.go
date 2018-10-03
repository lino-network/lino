// Package errors deals with all types of error encountered in the backend.
package errors

import (
	"fmt"
	"runtime"
)

// CodeType represents the type of the error.
type CodeType int

// IsOK returns true if the code type is OK.
func (code CodeType) IsOK() bool {
	return code == CodeOK
}

// Code types.
const (
	CodeOK CodeType = iota // 0
	CodeUnknown
	CodeAlreadyExists
	CodeFailedPrecondition
	CodeInternal // used by db operation
	CodeInvalidArgument
	CodeNotFound
	CodePermissionDenied
	CodeAPIError
	CodeUnavailable
	CodeUnablePrepareStatement // used by db sql prepare statement
	CodeFailedToSendEmail
	CodeFailedToSendCashout
	CodeInvalidInstantCashout
	CodeUserNotFound
	CodeFailedToScan
)

// NOTE: Don't stringer this, we'll put better messages in later.
func CodeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeUnknown:
		return "Unknown request"
	case CodeAlreadyExists:
		return "Alreay exists"
	case CodeFailedPrecondition:
		return "Failed precondition error"
	case CodeInternal:
		return "Internal error"
	case CodeInvalidArgument:
		return "Invalid argument"
	case CodeNotFound:
		return "Not found"
	case CodePermissionDenied:
		return "Permission denied"
	case CodeAPIError:
		return "API error"
	case CodeUnavailable:
		return "Unavailable"
	case CodeUnablePrepareStatement:
		return "Unable prepare sql statement"
	case CodeFailedToSendEmail:
		return "Failed to send email"
	case CodeFailedToSendCashout:
		return "Failed to send cashout"
	case CodeInvalidInstantCashout:
		return "Invalid instant cashout"
	default:
		return fmt.Sprintf("Unknown code %d", code)
	}
}

//--------------------------------------------------------------------------------
func AlreadyExists(msg string) Error {
	return newError(CodeAlreadyExists, msg)
}

func AlreadyExistsf(format string, args ...interface{}) Error {
	return newError(CodeAlreadyExists, fmt.Sprintf(format, args...))
}

func FailedPrecondition(msg string) Error {
	return newError(CodeFailedPrecondition, msg)
}

func FailedPreconditionf(format string, args ...interface{}) Error {
	return newError(CodeFailedPrecondition, fmt.Sprintf(format, args...))
}

func Internal(msg string) Error {
	return newError(CodeInternal, msg)
}

func Internalf(format string, args ...interface{}) Error {
	return newError(CodeInternal, fmt.Sprintf(format, args...))
}

func InvalidArgument(msg string) Error {
	return newError(CodeInvalidArgument, msg)
}

func InvalidArgumentf(format string, args ...interface{}) Error {
	return newError(CodeInvalidArgument, fmt.Sprintf(format, args...))
}

func NotFound(msg string) Error {
	return newError(CodeNotFound, msg)
}

func NotFoundf(format string, args ...interface{}) Error {
	return newError(CodeNotFound, fmt.Sprintf(format, args...))
}

func PermissionDenied(msg string) Error {
	return newError(CodePermissionDenied, msg)
}

func PermissionDeniedf(format string, args ...interface{}) Error {
	return newError(CodePermissionDenied, fmt.Sprintf(format, args...))
}

func APIError(msg string) Error {
	return newError(CodeAPIError, msg)
}

func APIErrorf(format string, args ...interface{}) Error {
	return newError(CodeAPIError, fmt.Sprintf(format, args...))
}

func Unavailable(msg string) Error {
	return newError(CodeUnavailable, msg)
}

func Unavailablef(format string, args ...interface{}) Error {
	return newError(CodeUnavailable, fmt.Sprintf(format, args...))
}

func UnablePrepareStatement(msg string) Error {
	return newError(CodeUnablePrepareStatement, msg)
}

func UnablePrepareStatementf(format string, args ...interface{}) Error {
	return newError(CodeUnablePrepareStatement, fmt.Sprintf(format, args...))
}

func FailedToSendEmail(msg string) Error {
	return newError(CodeFailedToSendEmail, msg)
}

func FailedToSendEmailf(format string, args ...interface{}) Error {
	return newError(CodeFailedToSendEmail, fmt.Sprintf(format, args...))
}

func FailedToSendCashout(msg string) Error {
	return newError(CodeFailedToSendCashout, msg)
}

func FailedToSendCashoutf(format string, args ...interface{}) Error {
	return newError(CodeFailedToSendCashout, fmt.Sprintf(format, args...))
}

func InvalidInstantCashout(msg string) Error {
	return newError(CodeInvalidInstantCashout, msg)
}

func InvalidInstantCashoutf(format string, args ...interface{}) Error {
	return newError(CodeInvalidInstantCashout, fmt.Sprintf(format, args...))
}

//----------------------------------------
// Error & serverError

// Error interface for all DLive errors
type Error interface {
	Error() string
	CodeType() CodeType
	Trace(msg string) Error
	Tracef(msg string, args ...interface{}) Error
	TraceCause(cause error, msg string) Error
	Cause() error
}

// NewError creates a new Error
func NewError(code CodeType, msg string) Error {
	return newError(code, msg)
}

// NewErrorf creates a new formatted Error
func NewErrorf(code CodeType, format string, args ...interface{}) Error {
	return newError(code, fmt.Sprintf(format, args...))
}

type traceItem struct {
	msg      string
	filename string
	lineno   int
}

func (ti traceItem) String() string {
	return fmt.Sprintf("%v:%v %v", ti.filename, ti.lineno, ti.msg)
}

// serverError is the customized Error used in the backend.
type serverError struct {
	code   CodeType
	msg    string
	cause  error
	traces []traceItem
}

func newError(code CodeType, msg string) *serverError {
	// TODO capture stacktrace if ENV is set.
	if msg == "" {
		msg = CodeToDefaultMsg(code)
	}
	return &serverError{
		code:   code,
		msg:    msg,
		cause:  nil,
		traces: nil,
	}
}

// Error returns error details.
func (err *serverError) Error() string {
	traceLog := ""
	for _, ti := range err.traces {
		traceLog += ti.String() + "\n"
	}
	return fmt.Sprintf("Error{%d:%s,%v\ntrace:\n%v}", err.code, err.msg, err.cause, traceLog)
}

// CodeType returns the code of error.
func (err *serverError) CodeType() CodeType {
	return err.code
}

// Trace adds tracing information with msg.
func (err *serverError) Trace(msg string) Error {
	return err.doTrace(msg, 2)
}

// Tracef adds tracing information with formatted msg.
func (err *serverError) Tracef(format string, arg ...interface{}) Error {
	msg := fmt.Sprintf(format, arg...)
	return err.doTrace(msg, 2)
}

// TraceCause adds tracing information with cause and msg.
func (err *serverError) TraceCause(cause error, msg string) Error {
	err.cause = cause
	return err.doTrace(msg, 2)
}

func (err *serverError) doTrace(msg string, depth int) Error {
	_, fileName, line, ok := runtime.Caller(depth)
	if !ok {
		if fileName == "" {
			fileName = "<unknown>"
		}
		if line <= 0 {
			line = -1
		}
	}
	// Do not include the whole stack trace.
	err.traces = append(err.traces, traceItem{
		filename: fileName,
		lineno:   line,
		msg:      msg,
	})
	return err
}

// Cause returns the cause of error.
func (err *serverError) Cause() error {
	return err.cause
}
