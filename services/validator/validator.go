package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	VALIDATOR_VALUE_SIGN       = "="
	VALIDATOR_RANGE_SPLIT_SIGN = ","
	VALIDATOR_PLACEHOLDER      = "_"
)

func NewValidator() *Validator {
	return &Validator{
		tagName:       "validator",
		splitSign:     "||",
		lazy:          false,
		validatorsMap: validatorsMap,
	}
}

type Validator struct {
	tagName       string
	splitSign     string
	lazy          bool
	validatorsMap map[string]IValidator
}

func (v *Validator) Validate(a interface{}) (errs []error) {
	errs = v.validate(a)
	return
}

// 添加验证器
func (v *Validator) AddValidator(key string, validator IValidator) error {
	if _, ok := v.validatorsMap[key]; ok {
		return fmt.Errorf(VALIDATOR_ALREADY_EXISTED, key)
	}
	v.validatorsMap[key] = validator
	return nil
}

// 设置惰性验证
func (v *Validator) SetLazy(lazy bool) {
	v.lazy = lazy
}

func (v *Validator) validate(a interface{}) (errs []error) {
	tp := reflect.TypeOf(a)
	vl := reflect.ValueOf(a)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		vl = vl.Elem()
	}

	switch tp.Kind() {
	case reflect.Struct:
		errs = v.handleStut(vl)
	}

	return
}

func (v *Validator) handleStut(vl reflect.Value) (errs []error) {
	numField := vl.NumField()
	if numField <= 0 {
		errs = append(errs, fmt.Errorf(STRUCT_EMPTY, vl.Type().Name()))
		return
	}

	for i := 0; i < numField; i++ {
		fieldTypeInfo := vl.Type().Field(i)
		tag := fieldTypeInfo.Tag.Get(v.tagName)

		fieldValue := vl.Field(i)
		// fieldTypeKind := fieldInfo.Type().Kind()

		if err := v.hanleVerifyFromTag(tag, fieldTypeInfo, fieldValue); err != nil {
			errs = append(errs, err...)
			if v.lazy {
				return
			}
			continue
		}
	}

	return
}

func (va *Validator) hanleVerifyFromTag(tag string, field reflect.StructField, value reflect.Value) (errs []error) {
	args := strings.Split(tag, va.splitSign)

	isRequired := false
	for _, v := range args {
		if v == "required" {
			isRequired = true
			break
		}
	}

	for _, v := range args {
		vKey := v
		vArgs := make([]string, 0)

		idx := strings.Index(v, VALIDATOR_VALUE_SIGN)
		if idx != -1 {
			vKey = v[0:idx]
			vArgs = strings.Split(v[idx+1:], VALIDATOR_RANGE_SPLIT_SIGN)
		}
		vali, ok := va.validatorsMap[vKey];
		if !ok {
			errs = append(errs, fmt.Errorf("[%s] validator not exist", vKey))
			if va.lazy {
				return
			}
			continue
		}

		if err := vali.Validate(field.Name, value, isRequired, vArgs...); err != nil {
			errs = append(errs, err)
			if va.lazy {
				return
			}
			continue
		}
	}
	return
}

func (v *Validator) checkIsMulti(value reflect.Value) bool {
}
