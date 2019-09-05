package validator

import "reflect"

type IValidator interface {
	Validate(field string, value reflect.Value, isRequired bool, args ...string) error
}

var validatorsMap = map[string]IValidator{
	"required": &RequiredValidator{},
	"string":   &StringValidator{},
	"number":   &NumberValidator{},
	"array":    &ArrayValidator{},
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
	Equal
}

func (sv *StringValidator) Validate(field string, value reflect.Value, isRequired bool, args ...string) error {
	eMsg := "[name] is not a string"
	if sv.EMsg != "" {
		eMsg = sv.EMsg
	}

	if value.Kind() != reflect.String {
		return formatError(eMsg, field)
	}

	argsLen := len(args)
	if argsLen == 1 {
		if !sv.StringEqual(value, args[0]) {
			return formatError("[name] is not equal to "+args[0], field)
		}
	}

	return nil
}

/*
NumberValidator
*/
type NumberValidator struct {
	EMsg string
	Equal
}

func (nv *NumberValidator) Validate(field string, value reflect.Value, isRequired bool, args ...string) error {
	eMsg := "[name] is not a number"
	if nv.EMsg != "" {
		eMsg = nv.EMsg
	}

	if !checkIsNumberKind(value.Kind()) {
		return formatError(eMsg, field)
	}

	argsLen := len(args)
	if argsLen == 1 {
		if !nv.NumberEqual(value, args[0]) {
			return formatError("[name] is not equal to "+args[0], field)
		}
	}

	return nil
}

/*
ArrayValidator
*/
type ArrayValidator struct {
	EMsg string
	Equal
}

func (av *ArrayValidator) Validate(field string, value reflect.Value, isRequired bool,
	args ...string) error {
	eMsg := "[name] is not a array/slice/map"
	if av.EMsg != "" {
		eMsg = av.EMsg
	}

	if !checkIsMultiKind(value.Kind()) {
		return formatError(eMsg, field)
	}

	return nil
}
