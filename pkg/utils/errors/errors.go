package utils

import (
	"errors"
	"fmt"
	"net/http"
)

type Error interface {
	Error() string
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call Unwrap on either.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

func IsItemNotFoundError(err error) bool {
	_, ok := err.(*ItemNotFoundError)
	return ok
}

func ParseStatusCode(err Error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err.(type) {
	case *BadRequestError:
		return http.StatusBadRequest
	case *ItemNotFoundError:
		return http.StatusNotFound
	case *InternalServerError:
		return http.StatusInternalServerError
	case *ForbiddenError:
		return http.StatusForbidden
	case *IsConflictError:
		return http.StatusConflict
	default:
		return http.StatusTeapot
	}
}

func ParseErrorType(err Error) string {
	if err == nil {
		return "No error"
	}
	switch err.(type) {
	case *BadRequestError:
		return "Bad request"
	case *ItemNotFoundError:
		return "Item not found error"
	case *InternalServerError:
		return "Internal server error"
	case *IsConflictError:
		return "Conflict error"
	case *ForbiddenError:
		return "Forbidden error"
	default:
		return "Bad request"
	}
}

func New(msg interface{}) error {
	if msg == nil {
		return nil
	}
	return fmt.Errorf("%v", msg)
}

type BadRequestError struct {
	msg string
}

type ItemNotFoundError struct {
	msg string
}

type InternalServerError struct {
	msg string
}

type ForbiddenError struct {
}

type IsConflictError struct {
	msg string
}

func NewBadRequestError(msg interface{}) Error {
	if msg == nil {
		return nil
	}
	return &BadRequestError{msg: fmt.Sprintf("%v", msg)}
}

func NewBadRequestErrorf(format string, vals ...interface{}) Error {
	return &BadRequestError{msg: fmt.Sprintf(format, vals...)}
}

func NewItemNotFoundError(msg interface{}) Error {
	if msg == nil {
		return nil
	}
	return &ItemNotFoundError{msg: fmt.Sprintf("%v", msg)}
}

func NewItemNotFoundErrorf(format string, vals ...interface{}) Error {
	return &ItemNotFoundError{msg: fmt.Sprintf(format, vals...)}
}

func NewInternalServerError(msg interface{}) Error {
	if msg == nil {
		return nil
	}
	return &InternalServerError{msg: fmt.Sprintf("%v", msg)}
}

func NewInternalServerErrorf(format string, vals ...interface{}) Error {
	return &InternalServerError{msg: fmt.Sprintf(format, vals...)}
}

func NewForbiddenError() Error {
	return &ForbiddenError{}
}

func NewIsConflictError(msg interface{}) Error {
	if msg == nil {
		return nil
	}
	return &IsConflictError{msg: fmt.Sprintf("%v", msg)}
}

func (e *BadRequestError) Error() string {
	return e.msg
}

func (e *ItemNotFoundError) Error() string {
	return e.msg
}

func (e *InternalServerError) Error() string {
	return e.msg
}

func (e *ForbiddenError) Error() string {
	return "forbidden"
}

func (e *IsConflictError) Error() string {
	return e.msg
}
