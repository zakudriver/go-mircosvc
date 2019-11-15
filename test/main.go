package main

import (
	"fmt"
	"time"
)

type A struct {
	a int
}

type B struct {
	a int
}

func main() {
	fmt.Println(time.Now().String())
}

func ss() *A {
	return &A{}
}
