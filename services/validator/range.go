package validator

import (
	"reflect"
	"strconv"
)

/****************************************************
 * range 验证错误提示 map
 ****************************************************/
var stringErrMap = map[string]string{
	"lessThan": "[name] should be less than [max] chars long",
	"equal":    "[name] should be equal [min] chars long",
	"atLeast":  "[name] should be at least [min] chars long",
	"between":  "[name] should be betwween [min] and [max] chars long",
}

var numberErrMap = map[string]string{
	"lessThan": "[name] should be less than [max]",
	"equal":    "[name] should be equal [min]",
	"atLeast":  "[name] should be at least [min]",
	"between":  "[name] should be betwween [min] and [max]",
}

var arrayErrMap = map[string]string{
	"lessThan": "array [name] length should be less than [max]",
	"equal":    "array [name] length should be equal [min]",
	"atLeast":  "array [name] length should be at least [min]",
	"between":  "array [name] length should be betwween [min] and [max]",
}

type Range struct {
	min    string
	max    string
	errMap map[string]string
}

func (r *Range) SetRangeIndex(field string, args ...string) {
	argsLen := len(args)
	if argsLen != 2 {
		panic("args length should be equal 2")
	}

	if args[0] != VALIDATOR_PLACEHOLDER && args[1] != VALIDATOR_PLACEHOLDER {
		min, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			panic("min must be int/float")
		}

		max, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			panic("max must be int/float")
		}

		if min >= max {
			panic("max must be greater than min")
		}
	}

	if args[0] == VALIDATOR_PLACEHOLDER && args[1] == VALIDATOR_PLACEHOLDER {
		panic(`min and max can't be "_" at the same time`)
	}

	r.min = args[0]
	r.max = args[1]
	r.errMap["min"] = args[0]
	r.errMap["max"] = args[1]
	r.errMap["name"] = field
}

func (r *Range) CompareFloat(value reflect.Value, field string) error {
	vl := value.Interface()
	var f float64
	if value.Kind() == reflect.Float32 {
		f = float64(vl.(float32))
	} else if value.Kind() == reflect.Float64 {
		f = vl.(float64)
	} else {
		return formatError("[name] not is float32/float64", field)
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
		status = 2
	}

	errKey := ""
	switch status {
	case 0:
		if min > f || max < f {
			errKey = "between"
		}
		break
	default:
		return nil
	}

	if errKey != "" {
		return formatMapError(numberErrMap[errKey], r.errMap)
	}

	return nil
}
