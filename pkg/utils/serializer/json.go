package serializer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Error interface {
	Error() string
}

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

type (
	BadRequestError struct {
		msg string
	}

	ItemNotFoundError struct {
		msg string
	}

	InternalServerError struct {
		msg string
	}

	ForbiddenError struct{}

	IsConflictError struct {
		msg string
	}
)

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

// swagger:model message
type JsonMessage struct {
	Message string `json:"message"`
}

func SendJSON(w http.ResponseWriter, status int, object interface{}) error {
	// set headers
	SetCorsHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// write response
	if err := json.NewEncoder(w).Encode(object); err != nil {
		return err
	}
	return nil
}

func SendNoContent(w http.ResponseWriter) error {
	// set headers
	SetCorsHeaders(w)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func SendError(w http.ResponseWriter, err Error) error {
	return SendJSON(w, ParseStatusCode(err),
		&map[string]string{"type": ParseErrorType(err), "message": err.Error()},
	)
}

func SetCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func ParseJsonBody(body io.ReadCloser, value interface{}) error {
	if err := json.NewDecoder(body).Decode(value); err != nil && !Is(err, io.EOF) {
		return NewBadRequestError(err)
	}
	return nil
}
