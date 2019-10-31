package main

import "fmt"

func main() {
	// ch := make(chan int, 1)
	type A struct {
		Uid string
		b   int
	}

	var a interface{} = A{Uid: "11"}
	fmt.Println(a.(A))
}
