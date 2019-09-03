package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func NewValidator() *Validator {
	return &Validator{
		tagName:     "validator",
		splitSymbol: "||",
		lazy:        false,
	}
}

type Validator struct {
	tagName     string
	splitSymbol string
	lazy        bool
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
		tag := fieldTypeInfo.Tag.Get(v.tagName)

		fieldInfo := vl.Field(i)
		fieldTypeKind := fieldInfo.Type().Kind()

		fmt.Println(fieldTypeInfo.Tag)
		fmt.Println(fieldInfo)
		fmt.Println(fieldTypeKind)
	}

	return nil
}

func (v *Validator) hanleVerifyFromTag(tag string, field reflect.StructField, value reflect.Value) {

	args := strings.Split(tag, v.splitSymbol)
	for _, v := range args {

		idx := strings.Index(v, "=")
		if idx!=-1{


		}
	}

}

func (v *Validator) checkArrayValueIsMulti(value reflect.Value) {
}
