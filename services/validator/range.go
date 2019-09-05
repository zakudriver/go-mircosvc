package validator

import (
	"errors"
	"reflect"
	"strconv"
)

/****************************************************
 * range 验证错误提示 map
 ****************************************************/
var stringErrorMap = map[string]string{
	"lessThan": "[name] should be less than [max] chars long",
	"equal":    "[name] should be equal [min] chars long",
	"atLeast":  "[name] should be at least [min] chars long",
	"between":  "[name] should be betwween [min] and [max] chars long",
}

var numberErrorMap = map[string]string{
	"lessThan": "[name] should be less than [max]",
	"equal":    "[name] should be equal [min]",
	"atLeast":  "[name] should be at least [min]",
	"between":  "[name] should be betwween [min] and [max]",
}

var arrayErrorMap = map[string]string{
	"lessThan": "array [name] length should be less than [max]",
	"equal":    "array [name] length should be equal [min]",
	"atLeast":  "array [name] length should be at least [min]",
	"between":  "array [name] length should be betwween [min] and [max]",
}

type Range struct {
	min string
	max string
}

func (r *Range) SetRangeIndex(field string, args ...string) {
	argsLen := len(args)
	if argsLen != 2 {
		panic("args length should be equal 2")
	}

	r.min = args[0]
	r.max = args[1]
}

func (r *Range) CompareFloat(value reflect.Value, field string) error {
	if r.min == VALIDATOR_PLACEHOLDER && r.max == VALIDATOR_PLACEHOLDER {
		return errors.New(`min and max can't be "_" at the same time`)
	}
	status := 0
	var min float64
	var max float64

	if r.min != VALIDATOR_PLACEHOLDER {
		f, err := strconv.ParseFloat(r.min, 64)
		if err != nil {
			return err
		}
		min = f
		status = 1
	}

	if r.max != VALIDATOR_PLACEHOLDER {
		f, err := strconv.ParseFloat(r.max, 64)
		if err != nil {
			return err
		}
		max = f
		status = 1
	}

	vl := value.Interface()
	var f float64
	if value.Kind() == reflect.Float32 {
		f = float64(vl.(float32))
	} else if value.Kind() == reflect.Float64 {
		f = vl.(float64)
	}

	return nil
}
