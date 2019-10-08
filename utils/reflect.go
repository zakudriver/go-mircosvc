package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	MAP_TAG_NAME = "map"
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
func Struct2Map(a interface{}) map[string]interface{} {
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		m[t.Field(i).Name] = v.Field(i).Interface()
	}
	return m
}

// 根据tag struct -> map
func Struct2MapFromTag(a interface{}) map[string]interface{} {
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	m := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if tag := t.Field(i).Tag.Get(MAP_TAG_NAME); tag != "" {
			m[tag] = v.Field(i).Interface()
		}
	}
	return m
}

// 根据struts解析环境变量
func ParseEnvForTag(a interface{}, tagName string) (err error) {
	tp := reflect.TypeOf(a)
	if tp.Kind() != reflect.Ptr && tp.Elem().Kind() != reflect.Struct {
		err = errors.New("[ParseEnvForTag] params must be *Struct")
		return
	}

	tp = tp.Elem()
	vl := reflect.ValueOf(a).Elem()

	for i := 0; i < vl.NumField(); i++ {
		fieldTypeInfo := vl.Type().Field(i)
		fieldValue := vl.Field(i)
		tag := fieldTypeInfo.Tag.Get(tagName)

		if tag != "" {
			if fieldValue.Kind() != reflect.String {
				// err = errors.New("[ParseEnvForTag] Struct property must be String")
				// return
				continue
			}
			args := strings.Split(tag, "=")
			if len(args) == 0 {
				continue
			}

			env := ""
			if len(args) == 1 {
				env = os.Getenv(args[0])
			} else {
				if v := os.Getenv(args[0]); v != "" {
					env = v
				} else {
					env = args[1]
				}
			}
			fieldValue.Set(reflect.ValueOf(env))
		}

	}

	return
}
