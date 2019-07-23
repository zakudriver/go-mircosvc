package utils

import (
	"strconv"
	"time"
)

// 字符串名称月份 -> 整数月份
func Str2IntForMonth(month string) int {
	var data = map[string]int{
		"January":   0,
		"February":  1,
		"March":     2,
		"April":     3,
		"May":       4,
		"June":      5,
		"July":      6,
		"August":    7,
		"September": 8,
		"October":   9,
		"November":  10,
		"December":  11,
	}
	return data[month]
}

// 获取当天 年-月-日 格式
func GetTodayYMD(sep string) string {
	now := time.Now()
	y := now.Year()
	m := Str2IntForMonth(now.Month().String())
	d := now.Day()

	var mStr string
	var dStr string
	if m < 9 {
		mStr = "0" + strconv.Itoa(m+1)
	} else {
		mStr = strconv.Itoa(m + 1)
	}

	if d < 10 {
		dStr = "0" + strconv.Itoa(d)
	} else {
		dStr = strconv.Itoa(d)
	}
	return strconv.Itoa(y) + sep + mStr + sep + dStr
}
