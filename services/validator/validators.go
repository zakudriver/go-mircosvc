package validator

import "reflect"

type IValidator interface {
	Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error)
}
