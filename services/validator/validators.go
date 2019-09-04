package validator

import "reflect"

type IValidator interface {
	Validate(field string, value reflect.Value, isRequired bool, args ...string) error
}

var validatorsMap = map[string]IValidator{
	"required": &RequiredValidator{},
	"string":   &StringValidator{},
	"number":   &NumberValidator{},
}

/*
RequiredValidator
*/
type RequiredValidator struct {
	EMsg string
}

func (r *RequiredValidator) Validate(field string, value reflect.Value, isRequired bool, args ...string) error {
	eMsg := "[name] is must required"
	if r.EMsg != "" {
		eMsg = r.EMsg
	}

	if checkIsZoreValue(value) {
		return formatError(eMsg, field)
	}

	return nil
}

/*
StringValidator
*/
type StringValidator struct {
	EMsg string
}

func (s *StringValidator) Validate(field string, value reflect.Value, isRequired bool, args ...string) error {
	eMsg := "[name] is not a string"
	if s.EMsg != "" {
		eMsg = s.EMsg
	}

	if value.Kind() != reflect.String {
		return formatError(eMsg, field)
	}

	return nil
}

/*
NumberValidator
*/
type NumberValidator struct {
	EMsg string
}

func (nv *NumberValidator) Validate(field string, value reflect.Value, isRequired bool, args ...string) error {
	eMsg := "[name] is not a number"
	if nv.EMsg != "" {
		eMsg = nv.EMsg
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		break
	default:
		return formatError(eMsg, field)
	}

	return nil
}
