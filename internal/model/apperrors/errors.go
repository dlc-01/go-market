package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type typeCode string

const (
	Unauthorized         typeCode = "UNAUTHORIZED"
	BadRequest           typeCode = "BAD_REQUEST"
	Conflict             typeCode = "CONFLICT"
	Internal             typeCode = "INTERNAL"
	NotFound             typeCode = "NOTFOUND"
	ServiceUnavailable   typeCode = "SERVICE_UNAVAILABLE"
	UnsupportedMediaType typeCode = "UNSUPPORTED_MEDIA_TYPE"
	UnprocessablContent  typeCode = "UNPROCESSABLE_CONTENT"
	StatusOK             typeCode = "STATUS_OK"
	Accepted             typeCode = "ACCEPTED"
	NoContent            typeCode = "NO_CONTENT"
	PaymentRequired      typeCode = "PAYMENT_REQUIRED"
)

type Error struct {
	TypCode typeCode `json:"TypCode"`
	Message string   `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.TypCode {
	case StatusOK:
		return http.StatusOK
	case Accepted:
		return http.StatusAccepted
	case NoContent:
		return http.StatusNoContent
	case PaymentRequired:
		return http.StatusPaymentRequired
	case Unauthorized:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case NotFound:
		return http.StatusNotFound
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case UnprocessablContent:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

func NewUnauthorized(reason string) *Error {
	return &Error{
		TypCode: Unauthorized,
		Message: reason,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		TypCode: BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reason),
	}
}

func NewConflict(value string) *Error {
	return &Error{
		TypCode: Conflict,
		Message: fmt.Sprintf("%v already exists", value),
	}
}

func NewInternal() *Error {
	return &Error{
		TypCode: Internal,
		Message: fmt.Sprint("Internal server error."),
	}
}

func NewNotFound(name string, value string) *Error {
	return &Error{
		TypCode: NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

func NewServiceUnavailable() *Error {
	return &Error{
		TypCode: ServiceUnavailable,
		Message: fmt.Sprint("Service unavailable or timed out"),
	}
}

func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		TypCode: UnsupportedMediaType,
		Message: reason,
	}
}
func NewUnprocessableContent(reason string) *Error {
	return &Error{
		TypCode: UnprocessablContent,
		Message: fmt.Sprintf("Unprocessable Content. Reasone: %v", reason),
	}
}
func NewStatusOK() *Error {
	return &Error{
		TypCode: StatusOK,
	}
}

func NewAccepted() *Error {
	return &Error{
		TypCode: Accepted,
	}
}

func NewNoContent() *Error {
	return &Error{
		TypCode: NoContent,
	}
}

func NewPaymentRequired(reason string) *Error {
	return &Error{
		TypCode: PaymentRequired,
		Message: fmt.Sprintf("Payment Required. Reasone: %v", reason),
	}
}
