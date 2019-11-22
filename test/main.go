package main

import (
	"fmt"
	"github.com/kum0/go-mircosvc/utils"
)

type A struct {
	a int
}

func (a *A) ToString() string {
	return ""
}

func (a *A) ToInt() string {
	return ""
}

type B struct {
	a int
}

type AA interface {
	ToString() string
	// BB
}

type BB interface {
	ToInt() string
}

func main() {
	a, _ := utils.NewUUID()
	fmt.Println(a)
}

func test(b BB) {
	b.ToInt()
}
