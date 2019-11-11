package main

import "fmt"

type A struct {
	a int
}

func main() {
	var a A
	fmt.Println(a)
}

func ss() *A {
	return &A{}
}
