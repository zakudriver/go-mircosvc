package validator

import (
	"reflect"
)

type Equal struct {
}

func (e *Equal) StringEqual(value reflect.Value, strVl string) bool {
	if value.Kind() != reflect.String {
		return false
	}
	return value.Interface() == strVl
}

func (e *Equal) NumberEqual(value reflect.Value, strVl string) bool {
	if !checkIsNumberKind(value.Kind()) {
		return false
	}

	return numberToString(value) == strVl
}
