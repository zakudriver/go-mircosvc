package validator

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
