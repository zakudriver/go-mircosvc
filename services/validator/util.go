package validator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Kind uint

const (
	TYPE_INVALID              = "[name] type invalid"
	PARAMS_INVALID            = "[name] params invalid"
	STRUCT_EMPTY              = "struct %v is empty"
	VALIDATOR_ALREADY_EXISTED = "[%s] validator already existed"
	// ERROR_NAME_PLACEHOLDER    = "name"

	Int_Kind Kind = iota
	Int8_Kind
	Int16_Kind
	Int32_Kind
	Int64_Kind
	Uint_Kind
	Uint8_Kind
	Uint16_Kind
	Uint32_Kind
	Uint64_Kind
	Float32_Kind
	Float64_Kind
	Array_Kind
	Slice_Kind
	Map_Kind
)

var (
	numberKindMap = map[reflect.Kind]Kind{
		reflect.Int:     Int_Kind,
		reflect.Int8:    Int8_Kind,
		reflect.Int16:   Int16_Kind,
		reflect.Int32:   Int32_Kind,
		reflect.Int64:   Int64_Kind,
		reflect.Uint:    Uint_Kind,
		reflect.Uint8:   Uint8_Kind,
		reflect.Uint16:  Uint16_Kind,
		reflect.Uint32:  Uint32_Kind,
		reflect.Uint64:  Uint64_Kind,
		reflect.Float32: Float32_Kind,
		reflect.Float64: Float64_Kind,
	}

	multiKindMap = map[reflect.Kind]Kind{
		reflect.Array: Array_Kind,
		reflect.Slice: Slice_Kind,
		reflect.Map:   Map_Kind,
	}

	numberKindToStringFunc = map[reflect.Kind]func(a interface{}) string{
		reflect.Int: func(a interface{}) string {
			return strconv.Itoa(a.(int))
		},
		reflect.Int8: func(a interface{}) string {
			return strconv.Itoa(int(a.(int8)))
		},
		reflect.Int16: func(a interface{}) string {
			return strconv.Itoa(int(a.(int16)))
		},
		reflect.Int32: func(a interface{}) string {
			return strconv.Itoa(int(a.(int32)))
		},
		reflect.Int64: func(a interface{}) string {
			return strconv.Itoa(int(a.(int64)))
		},
		reflect.Uint: func(a interface{}) string {
			return strconv.Itoa(int(a.(uint)))
		},
		reflect.Uint8: func(a interface{}) string {
			return strconv.Itoa(int(a.(uint8)))
		},
		reflect.Uint16: func(a interface{}) string {
			return strconv.Itoa(int(a.(uint16)))
		},
		reflect.Uint32: func(a interface{}) string {
			return strconv.Itoa(int(a.(uint32)))
		},
		reflect.Uint64: func(a interface{}) string {
			return strconv.Itoa(int(a.(uint64)))
		},
		reflect.Float32: func(a interface{}) string {
			return strconv.FormatFloat(float64(a.(float32)), 'f', -1, 64)
		},
		reflect.Float64: func(a interface{}) string {
			return strconv.FormatFloat(a.(float64), 'f', -1, 64)

		},
	}
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
	e := strings.Replace(format, "name", field, 1)
	return errors.New(e)
}

func formatMapError(format string, fieldMaps ...map[string]string) error {
	var params []string
	for _, mv := range fieldMaps {
		for k, v := range mv {
			if strings.Index(format, k) >= 0 {
				params = append(params, k, v)
			}
		}
	}
	replacer := strings.NewReplacer(params...)
	return errors.New(replacer.Replace(format))
}

// 判断 Array/Slice/Map
func checkIsMultiKind(kind reflect.Kind) bool {
	_, ok := multiKindMap[kind]
	return ok
}

// 判断 Int/Int8/Int16/Int32/Int64/Uint/Uint8/Uint16/Uint32/Uint64/Float32/Float64/
func checkIsNumberKind(kind reflect.Kind) bool {
	_, ok := numberKindMap[kind]
	return ok
}

func numberToString(value reflect.Value) string {
	if !checkIsNumberKind(value.Kind()) {
		return ""
	}
	if fn, ok := numberKindToStringFunc[value.Kind()]; ok {
		return fn(value.Interface())
	}

	return ""
}
