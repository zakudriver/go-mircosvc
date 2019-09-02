package main

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type T struct {
	Name string
	Age  int
}

func main() {
	a := T{"zz", 11}
	b:=[]string{"1"}

	ty := reflect.TypeOf(b)

	va := reflect.ValueOf(a)
	fmt.Println(ty)
	fmt.Println(va)

	k := ty.Kind()
	fmt.Println(k )
}

func handle() error {
	fmt.Println("start")
	ch := make(chan error)

	go func() {
		time.Sleep(time.Second * 2)
		ch <- errors.New("err2")
	}()

	go func() {
		time.Sleep(time.Second * 3)
		ch <- errors.New("err3")
	}()

	n := 2
	for range ch {
		n--
		if n == 0 {
			close(ch)
			// fmt.Println(0)
		}
		// if v != nil {
		// 	close(ch)
		// 	return v
		// }
	}

	fmt.Println("over")
	return nil
}
