package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func NewRand(size int) int {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < size; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	rnum, _ := strconv.Atoi(sb.String())
	return rnum
}
