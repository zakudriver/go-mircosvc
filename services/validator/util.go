package validator

import (
	"errors"
	"reflect"
	"strings"
)

const (
	STRUCT_EMPTY              = "struct %v is empty"
	VALIDATOR_ALREADY_EXISTED = "[%s] validator already existed"
	ERROR_NAME_PLACEHOLDER    = "name"
)

func checkIsZoreValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Array:
		return value.Len() == 0
	case reflect.Slice, reflect.Map:
		return value.Len() == 0 || value.IsNil()
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	default:
		panic("handleIsZoreValue func error")
	}
}

func formatError(format, field string) error {
	e := strings.Replace(format, ERROR_NAME_PLACEHOLDER, field, 1)
	return errors.New(e)
}

