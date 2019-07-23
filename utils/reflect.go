package utils

import (
	"fmt"
	"reflect"
)

func setField(target interface{}, k string, v interface{}) error {
	structData := reflect.ValueOf(target).Elem()
	fieldValue := structData.FieldByName(k)

	if !fieldValue.IsValid() {
		return fmt.Errorf("utils.setField() No such field: %s in %s ", k, reflect.TypeOf(target))
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field v ", k)
	}

	fieldType := fieldValue.Type()
	val := reflect.ValueOf(v)

	valTypeStr := val.Type().String()
	fieldTypeStr := fieldType.String()
	if valTypeStr == "float64" && fieldTypeStr == "int" {
		val = val.Convert(fieldType)
	} else if fieldType != val.Type() {
		return fmt.Errorf("Provided v type " + valTypeStr + " didn't match target field type " + fieldTypeStr)
	}
	fieldValue.Set(val)
	return nil
}

// json映射 -> struct
func JSON2Struct(m map[interface{}]interface{}, target interface{}) error {
	for k, v := range m {
		if err := setField(target, k.(string), v); err != nil {
			return err
		}
	}
	return nil
}

// struct -> map
func Struct2Map(stut interface{}) (r map[string]interface{}) {
	t := reflect.TypeOf(stut)
	v := reflect.ValueOf(stut)

	r = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		r[t.Field(i).Name] = v.Field(i).Interface()
	}
	return
}

// struct tag
func HandleStructTag(stut interface{}, tagName string) (tags []string) {
	t := reflect.TypeOf(stut)

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(tagName)
		tags = append(tags, tag)
	}

	return tags
}
