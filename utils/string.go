package utils

import "strconv"

func String2Int(s string) int {
	r, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		panic(err)
	}
	return int(r)
}
