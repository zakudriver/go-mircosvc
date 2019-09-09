package validator

import (
	"reflect"
	"strconv"
	"strings"
)

// range 验证错误提示 map

var (
	stringErrMap = map[string]string{
		"lessThan": "[name] should be less than [max] chars long",
		"greaterThan":  "[name] should be greater than [min] chars long",
		"between":  "[name] should be between [min] and [max] chars long",
	}

	numberErrMap = map[string]string{
		"lessThan":    "[name] should be less than [max]",
		"greaterThan": "[name] should be greater than [min]",
		"between":     "[name] should be between [min] and [max]",
	}

	multiErrMap = map[string]string{
		"lessThan": "[name] length should be less than [max]",
		"greaterThan":  "[name] length should be greater than [min]",
		"between":  "[name] length should be between [min] and [max]",
	}
)

type min struct {
	value     string
	isInclude bool
}

type max struct {
	value     string
	isInclude bool
}

type Range struct {
	max
	min
	field  string
	errMap map[string]string
}

func (r *Range) InitRange(field string, args ...string) {
	r.field = field

	argsLen := len(args)
	if argsLen != 2 {
		panic("args length should be equal 2")
	}

	minArr := make([]string, 0)
	if strings.Index(args[0], "[") == 0 {
		minArr = append(minArr, "[", args[0][1:])
	} else {
		minArr = append(minArr, "", args[0])
	}

	maxArr := make([]string, 0)
	if strings.Index(args[1], "]") == len(args[1])-1 {
		maxArr = append(maxArr, args[1][:len(args[1])-1], "]")
	} else {
		maxArr = append(maxArr, args[1], "")
	}

	if minArr[1] != VALIDATOR_PLACEHOLDER_SIGN && maxArr[0] != VALIDATOR_PLACEHOLDER_SIGN {
		min, err := strconv.ParseFloat(minArr[1], 64)
		if err != nil {
			panic("min must be int/float.")
		}

		max, err := strconv.ParseFloat(maxArr[0], 64)
		if err != nil {
			panic("max must be int/float.")
		}

		if min >= max {
			panic("max must be greater than min.")
		}
	}

	if minArr[1] == VALIDATOR_PLACEHOLDER_SIGN && maxArr[0] == VALIDATOR_PLACEHOLDER_SIGN {
		panic(`min and max can't be "_" at the same time`)
	}

	r.min.value = minArr[1]
	r.min.isInclude = minArr[0] == "["

	r.max.value = maxArr[0]
	r.max.isInclude = maxArr[1] == "]"

	r.errMap = make(map[string]string)
	r.errMap["min"] = minArr[1]
	r.errMap["max"] = maxArr[0]
	r.errMap["name"] = field
}

// 区间比较
func (r *Range) compareSize(f float64, errMap map[string]string) error {
	// status:0 -> min!="_" && max!="_"
	// status:1 -> min!="_" && max=="_"
	// status:2 -> min=="_" && max!="_"
	status := 0
	var min float64
	var max float64

	if r.min.value != VALIDATOR_PLACEHOLDER_SIGN {
		f, err := strconv.ParseFloat(r.min.value, 64)
		if err != nil {
			return err
		}
		min = f
		status = 1
	}

	if r.max.value != VALIDATOR_PLACEHOLDER_SIGN {
		f, err := strconv.ParseFloat(r.max.value, 64)
		if err != nil {
			return err
		}
		max = f
		status = 2
	}

	if r.min.value != VALIDATOR_PLACEHOLDER_SIGN && r.max.value != VALIDATOR_PLACEHOLDER_SIGN {
		status = 0
	}

	minCondition := false
	if r.min.isInclude {
		minCondition = min > f
	} else {
		minCondition = min >= f
	}

	maxCondition := false
	if r.max.isInclude {
		maxCondition = max < f
	} else {
		maxCondition = min <= f
	}

	errKey := ""
	switch status {
	case 0:
		if minCondition || maxCondition {
			errKey = "between"
		}
		break
	case 1:
		if minCondition {
			errKey = "greaterThan"
		}
		break
	case 2:
		if maxCondition {
			errKey = "lessThan"
		}
		break
	}

	if errKey != "" {
		return formatMapError(errMap[errKey], r.errMap)
	}

	return nil
}

func (r *Range) CompareNumberRange(value reflect.Value) (err error) {
	if !checkIsNumberKind(value.Kind()) {
		return formatError("[name] not is int/int8/int16/int32/int64/uint/uint8/uint16/uint32/uint64/float32/float64", r.field)
	}

	vl := value.Interface()
	var f float64
	if value.Kind() == reflect.Float32 {
		f = float64(vl.(float32))
	} else if value.Kind() == reflect.Float64 {
		f = vl.(float64)
	} else {
		v := numberToString(value)
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
	}

	return r.compareSize(f, numberErrMap)
}

func (r *Range) CompareStringRange(value reflect.Value) error {
	if value.Kind() != reflect.String {
		return formatError("[name] is not a string", r.field)
	}

	l := value.Len()
	return r.compareSize(float64(l), stringErrMap)
}

func (r *Range) CompareMultiRange(value reflect.Value) error {
	if !checkIsMultiKind(value.Kind()) {
		return formatError("[name] is not a array/slice/map", r.field)
	}

	len := value.Len()
	return r.compareSize(float64(len), multiErrMap)
}
