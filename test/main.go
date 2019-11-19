package main

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
	var a = &A{}
	test(a)
}

func test(b BB) {
	b.ToInt()
}
