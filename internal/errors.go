package internal

import "fmt"

type Error struct {
	orig error
	msg  string
	code ErrorCode
}

type ErrorCode uint

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeNotFound
	ErrorCodeInvalidArgument
)

func WrapErrorf(orig error, format string, code ErrorCode, args ...interface{}) error {
	return &Error{
		orig: orig,
		msg:  fmt.Sprintf(format, args...),
		code: code,
	}
}

func NewErrorf(format string, code ErrorCode, args ...interface{}) error {
	return WrapErrorf(nil, format, code, args...)
}

func (e *Error) Error() string {
	if e.orig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.orig)
	}

	return e.msg
}

func (e *Error) Unwrap() error {
	return e.orig
}

func (e *Error) Code() ErrorCode {
	return e.code
}
