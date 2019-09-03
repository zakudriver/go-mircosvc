package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Zhan9Yunhua/blog-svr/services/validator"
)

type T struct {
	Name string
	Age  int
}

type P struct {
	Name string `validator:"required||string=1,5"`
	Age  int    `validator:"required||integer=1,5"`
}

func main() {
	// a := T{"zz", 11}
	//
	// tp := reflect.TypeOf(&a)
	// vl := reflect.ValueOf(&a)
	//
	// fmt.Println(vl.Type().Elem())
	//
	// info := vl.Elem().Field(0)
	// field := tp.Elem().Field(0)
	// fmt.Println(info)
	// fmt.Printf("%+v\n", field.Name)

	tar := P{"zzz", 11,}
	vali := validator.NewValidator()

	vali.Validate(tar)
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
