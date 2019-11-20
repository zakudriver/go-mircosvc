package common

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrRouteArgs           = errors.New("error route arguments")
	ErrArgs                = errors.New("error arguments")
	ErrEmpty               = errors.New("empty input")
	ErrInterfaceConversion = errors.New("interface conversion error")
)

func decodeError(err error) (int, error) {
	s := err.Error()
	c := s[:3]
	i, err := strconv.Atoi(c)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Internal Server Error")
	}

	return i, errors.New(s[4:])
}

func NewError(code int, a interface{}) error {
	var text string

	switch v := a.(type) {
	case string:
		text = v
	case error:
		text = v.Error()
	default:
		text = "Throw Error"
	}

	if code < 100 || code > 999 {
		code = 500
	}

	return &Error{code, text}
}

type Error struct {
	code int
	text string
}

func (e *Error) Error() string {
	return strconv.Itoa(e.code) + " " + e.text
}
