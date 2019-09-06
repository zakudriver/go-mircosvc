package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type IValidator interface {
	Validate(field string, value reflect.Value, isRequired bool, args ...string) error
}

var validatorsMap = map[string]IValidator{
	"required": &RequiredValidator{},
	"string":   &StringValidator{},
	"number":   &NumberValidator{},
	"multi":    &MultiValidator{},
	"in":       &InValidator{},
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
	Range
}

func (sv *StringValidator) Validate(field string, value reflect.Value, isRequired bool, args ...string) error {
	eMsg := "[name] is not a string"
	if sv.EMsg != "" {
		eMsg = sv.EMsg
	}

	if value.Kind() != reflect.String {
		return formatError(eMsg, field)
	}

	if !isRequired && checkIsZoreValue(value) {
		return nil
	}

	argsLen := len(args)
	if argsLen == 1 {
		return sv.ValueEqual(field, value, args[0])
	}

	if argsLen == 2 {
		sv.InitRange(field, args...)
		return sv.CompareStringRange(value)
	}

	return nil
}

/*
NumberValidator
*/
type NumberValidator struct {
	EMsg string
	Equal
	Range
}

func (nv *NumberValidator) Validate(field string, value reflect.Value, isRequired bool, args ...string) error {
	eMsg := "[name] is not a number"
	if nv.EMsg != "" {
		eMsg = nv.EMsg
	}

	if !checkIsNumberKind(value.Kind()) {
		return formatError(eMsg, field)
	}

	if !isRequired && checkIsZoreValue(value) {
		return nil
	}

	argsLen := len(args)
	if argsLen == 1 {
		return nv.ValueEqual(field, value, args[0])
	}

	if argsLen == 2 {
		nv.InitRange(field, args...)
		return nv.CompareNumberRange(value)
	}

	return nil
}

/*
MultiValidator
*/
type MultiValidator struct {
	EMsg string
	Equal
	Range
}

func (av *MultiValidator) Validate(field string, value reflect.Value, isRequired bool,
	args ...string) error {
	eMsg := "[name] is not a array/slice/map"
	if av.EMsg != "" {
		eMsg = av.EMsg
	}

	if !checkIsMultiKind(value.Kind()) {
		return formatError(eMsg, field)
	}

	if !isRequired && checkIsZoreValue(value) {
		return nil
	}

	if len(args) == 1 {
		return av.MultiEqual(field, value, args...)
	}

	if len(args) == 2 {
		av.InitRange(field, args...)
		return av.CompareMultiRange(value)
	}

	return nil
}

/*
BoolValidator
*/
type BoolValidator struct {
	EMsg string
	Equal
}

func (bv *BoolValidator) Validate(field string, value reflect.Value, isRequired bool,
	args ...string) error {
	eMsg := "[name] is not a bool"
	if bv.EMsg != "" {
		eMsg = bv.EMsg
	}

	if value.Kind() != reflect.Bool {
		return formatError(eMsg, field)
	}

	if !isRequired && checkIsZoreValue(value) {
		return nil
	}

	if len(args) == 1 {
		return bv.ValueEqual(field, value, args[0])
	}

	return nil
}

/*
	InValidator

	仅支持 string,float,int,bool 类型
	或值类型为 string,float,int,bool 类型的array,slice,map
*/
type InValidator struct {
	EMsg string
	Equal
}

func (iv *InValidator) Validate(field string, value reflect.Value, isRequired bool,
	args ...string) (err error) {
	eMsg := "[name] is not in [value]"
	if iv.EMsg != "" {
		eMsg = iv.EMsg
	}

	if len(args) == 0 {
		return errors.New("[InValidator] validator must have params")
	}

	if !isRequired && checkIsZoreValue(value) {
		return nil
	}

	isIn := false

	switch kind := value.Kind(); {
	case !checkIsMultiKind(kind):
		s := fmt.Sprint(value.Interface())
		if strings.Index(args[0], s) >= 0 {
			isIn = true
		}
		break
	case kind == reflect.Array || kind == reflect.Slice:
		s := fmt.Sprint(value.Interface())
		argsArr := strings.Split(args[0], ",")
		for _, v := range argsArr {
			if strings.Index(s, v) >= 0 {
				isIn = true
				break
			}
		}
		break
	case kind == reflect.Map:
		ks := value.MapKeys()
		for _, v := range ks {
			s := fmt.Sprint(value.MapIndex(v))
			if strings.Index(args[0], s) >= 0 {
				isIn = true
				break
			}
		}
		break
	}

	if !isIn {
		return formatMapError(eMsg, map[string]string{"name": field, "value": args[0]})
	}

	return nil
}
