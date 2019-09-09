package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	VALIDATOR_ASSIGNMENT_SIGN   = "="
	VALIDATOR_PARAMS_SPLIT_SIGN = "|"
	VALIDATOR_PLACEHOLDER_SIGN  = "_"
	VALIDATOR_SPLIT_SIGN        = "||"
)

func NewValidator() *Validator {
	return &Validator{
		tagName:       "validator",
		lazy:          false,
		validatorsMap: validatorsMap,
	}
}

type Validator struct {
	tagName       string
	lazy          bool
	validatorsMap map[string]IValidator
}

func (v *Validator) Validate(a interface{}) (errs []error) {
	errs = v.validate(a)
	return
}

func (v *Validator) LazyValidate(a interface{}) (errs []error) {
	oldLazy := v.lazy
	v.lazy = true
	errs = v.validate(a)
	v.lazy = oldLazy
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

func (v *Validator) validate(a interface{}) (errs []error) {
	tp := reflect.TypeOf(a)
	vl := reflect.ValueOf(a)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
		vl = vl.Elem()
	}

	switch kind := tp.Kind(); {
	case reflect.Struct == kind:
		if reErrs := v.handleStruct(vl); len(reErrs) > 0 {
			errs = append(errs, reErrs...)
		}
	case checkIsMultiKind(kind):
		if reErrs := v.handleMulti(vl); len(reErrs) > 0 {
			errs = append(errs, reErrs...)
		}
	}

	return
}

func (v *Validator) handleStruct(value reflect.Value) (errs []error) {
	numField := value.NumField()
	if numField <= 0 {
		errs = append(errs, fmt.Errorf(STRUCT_EMPTY, value.Type().Name()))
		return
	}

	for i := 0; i < numField; i++ {
		fieldTypeInfo := value.Type().Field(i)
		fieldValue := value.Field(i)
		tag := fieldTypeInfo.Tag.Get(v.tagName)
		if tag != "" {
			// fieldTypeKind := fieldInfo.Type().Kind()

			if reErrs := v.hanleVerifyFromTag(tag, fieldTypeInfo, fieldValue); len(reErrs) > 0 {
				errs = append(errs, reErrs...)
				if v.lazy {
					return
				}
				continue
			}
		}

		if value.Kind() == reflect.Struct {
			if reErrs := v.validate(fieldValue.Interface()); len(reErrs) > 0 {
				errs = append(errs, reErrs...)
				if v.lazy {
					return
				}
			}
			continue
		}

		if reErrs := v.handleMulti(value); len(reErrs) > 0 {
			errs = append(errs, reErrs...)
		}
	}

	return
}

func (v *Validator) handleMulti(value reflect.Value) (errs []error) {
	if v.checkIsMulti(value) {
		for i := 0; i < value.Len(); i++ {
			if reErrs := v.validate(value.Index(i).Interface()); len(reErrs) > 0 {
				errs = append(errs, reErrs...)
				if v.lazy {
					return
				}
			}
		}
	}

	return
}

func (va *Validator) hanleVerifyFromTag(tag string, field reflect.StructField, value reflect.Value) (errs []error) {
	args := strings.Split(tag, VALIDATOR_SPLIT_SIGN)

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

		idx := strings.Index(v, VALIDATOR_ASSIGNMENT_SIGN)
		if idx != -1 {
			vKey = v[0:idx]
			vArgs = strings.Split(v[idx+1:], VALIDATOR_PARAMS_SPLIT_SIGN)
		}
		vali, ok := va.validatorsMap[vKey]
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

// 判断 arrar/slice/map的item是否是array/slice/map/struct
func (v *Validator) checkIsMulti(value reflect.Value) (ok bool) {
	if ok = checkIsMultiKind(value.Kind()); !ok {
		return
	}

	valueKind := value.Type().Elem().Kind()

	if ok = checkIsMultiKind(valueKind) || valueKind == reflect.Struct; !ok {
		return
	}

	return
}
