package crashy

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const (
	maxStackLength = 50
)

type (
	Error struct {
		code       string
		message    string
		stacktrace string
		origin     error
	}
	CodeMapper interface {
		Code() string
		Message() string
	}
	Wrapped interface {
		Unwrap() error
	}
	StackTracer interface {
		StackTrace() string
	}
)

func New(code string, message string) *Error {
	return &Error{
		code:       code,
		message:    message,
		origin:     errors.New(message),
		stacktrace: getStackTrace(2),
	}
}

func (e Error) Error() string {
	if e.code == "" {
		e.code = ErrCodeUnexpected
	}
	if e.message == "" {
		e.message = Message(e.code)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e Error) Unwrap() error {
	return e.origin
}

func (e Error) StackTrace() string {
	return e.stacktrace
}

func (e Error) Code() string {
	if e.code == "" {
		e.code = ErrCodeUnexpected
	}
	return e.code
}
func (e Error) Message() string {
	if e.code == "" {
		e.code = ErrCodeUnexpected
	}
	if e.message == "" {
		e.message = Message(e.code)
	}
	return e.message
}

func Is(err error, code string) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(CodeMapper); ok {
		return e.Code() == code
	}
	return false
}

func Wrap(err error, code string, message string) error {
	if err == nil {
		return New(code, message)
	}

	if e, ok := err.(Wrapped); ok { //we do not support multi wrapper
		err = e.Unwrap()
	}

	var st string
	if e, ok := err.(StackTracer); ok { //we keep original stacktrace
		st = e.StackTrace()
	}
	if st == "" {
		st = getStackTrace(2)
	}

	return &Error{
		code:       code,
		message:    message,
		origin:     err,
		stacktrace: st,
	}
}

func WithCode(err error, code string) error {
	return Wrap(err, code, "")
}

func WithFormat(err error, code string, format string, params ...interface{}) error {
	return Wrap(err, code, fmt.Sprintf(format, params...))
}

func getStackTrace(skip int) string {
	stackBuf := make([]uintptr, maxStackLength)
	length := runtime.Callers(skip, stackBuf[:])
	stack := stackBuf[:length]

	trace := ""
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = trace + fmt.Sprintf("\n\tFile: %s, Line: %d. Function: %s", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}
	return trace
}
