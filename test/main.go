package main

import (
	"fmt"

	"github.com/Zhan9Yunhua/blog-svr/utils"
)

func main() {
	// ch := make(chan int, 1)
	type A struct {
		Uid string
		b   int
	}

	type B struct {
		Uid string
	}

	a := A{"aaa", 11}
	b := &B{}
	fmt.Println(utils.StructCopy(a, 11))
	fmt.Println(b)
}
