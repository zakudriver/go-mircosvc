package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	EQUAL_ERR = "[name] should be equal [value]"
)

type Equal struct {
}

// 判断 String/Bool/Int/Int8/Int16/Int32/Int64/Uint/Uint8/Uint16/Uint32/Uint64/Float32/Float64
func (e *Equal) ValueEqual(field string, value reflect.Value, vl string) error {
	if fmt.Sprint(value.Interface()) != vl {
		return formatMapError(EQUAL_ERR, map[string]string{"name": field, "value": vl})
	}
	return nil
}

func (e *Equal) MultiEqual(field string, value reflect.Value, args ...string) error {
	valueKind := value.Type().Elem().Kind()
	if checkIsMultiKind(valueKind) {
		return formatError(TYPE_INVALID, field)
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		v := strings.Join(args, " ")
		if fmt.Sprint(value.Interface()) != v {
			return formatMapError(EQUAL_ERR, map[string]string{"name": field, "value": v})
		}
		break
	case reflect.Map:
		l := len(args)
		if strings.Index(args[0], "{") == 0 && strings.Index(args[l-1], "}") == len(args[l-1])-1 {
			args[0] = args[0][1:]
			args[l-1] = args[l-1][:len(args[l-1])-1]
			m := make(map[string]string)
			for _, v := range args {
				kv := strings.Split(v, ":")
				m[kv[0]] = kv[1]
			}

			for _, v := range value.MapKeys() {
				k := fmt.Sprint(v)
				if mv, ok := m[k]; ok {
					if mv == fmt.Sprint(value.MapIndex(v).Interface()) {
						continue
					}
				}
				return formatMapError(EQUAL_ERR, map[string]string{"name": field, "value": strings.Join(args, ",")})
			}

		} else {
			return formatError(PARAMS_INVALID, field)
		}
		break
	}

	return formatError(TYPE_INVALID, field)
}
