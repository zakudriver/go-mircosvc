package main

import (
	"fmt"
	"time"
)

type A struct {
	a int
}

func main() {
	fmt.Println(int(time.Now().Unix()))
}

func ss() *A {
	return &A{}
}
