package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	VALIDATOR_VALUE_SIGN       = "="
	VALIDATOR_RANGE_SPLIT_SIGN = ","
	VALIDATOR_IGNORE_SIGN      = "_"
)

func NewValidator() *Validator {
	return &Validator{
		tagName:       "validator",
		splitSign:     "||",
		lazy:          false,
		validatorsMap: make(map[string]IValidator),
	}
}

type Validator struct {
	tagName       string
	splitSign     string
	lazy          bool
	validatorsMap map[string]IValidator
}

func (v *Validator) Validate(a interface{}) {
	v.validate(a)
}

func (v *Validator) validate(a interface{}) {
	tp := reflect.TypeOf(a)
	vl := reflect.ValueOf(a)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		vl = vl.Elem()
	}

	switch tp.Kind() {
	case reflect.Struct:
		v.handleStut(vl)
	}
}

func (v *Validator) handleStut(vl reflect.Value) (errs []error) {
	numField := vl.NumField()
	if numField <= 0 {
		errs = append(errs, fmt.Errorf(STRUCT_EMPTY, vl.Type().Name()))
		return
	}

	for i := 0; i < numField; i++ {
		fieldTypeInfo := vl.Type().Field(i)
		// tag := fieldTypeInfo.Tag.Get(v.tagName)

		fieldInfo := vl.Field(i)
		fieldTypeKind := fieldInfo.Type().Kind()

		fmt.Println(fieldTypeInfo.Tag)
		fmt.Println(fieldInfo)
		fmt.Println(fieldTypeKind)
	}

	return nil
}

func (v *Validator) hanleVerifyFromTag(tag string, field reflect.StructField, value reflect.Value) {

	args := strings.Split(tag, v.splitSign)
	for _, v := range args {
		vKey := v
		vArgs := make([]string, 0)

		idx := strings.Index(v, VALIDATOR_VALUE_SIGN)
		if idx != -1 {
			vKey = v[0:idx]
			vArgs = strings.Split(v[idx+1:], VALIDATOR_RANGE_SPLIT_SIGN)
		}
	}

}

func (v *Validator) checkArrayValueIsMulti(value reflect.Value) {
}
