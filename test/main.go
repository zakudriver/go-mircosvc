package main

type A struct {
	a int
}

type B struct {
	a int
}

func main() {
	a := A{a: 1}
	b := B{}

	a = b
}

func ss() *A {
	return &A{}
}
