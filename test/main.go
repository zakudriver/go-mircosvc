package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type T struct {
	TName string `validator:"number"`
	TAge  int    `validator:"number"`
}

type P struct {
	Name string  `validator:"required"`
	Age  float64 `validator:"number=11.1"`
	Sex  int     `validator:"string"`
}

func main() {
	// a := P{Name: "zz", Age: 11.1,}
	//
	// vali := validator.NewValidator()
	//
	// err := vali.Validate(a)
	//
	// fmt.Printf("%+v\n", err)
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

func formatMapError(format string, fieldMap map[string]string) error {
	var params []string
	for k, v := range fieldMap {
		if strings.Index(format, k) >= 0 {
			params = append(params, k, v)
		}
	}
	replacer := strings.NewReplacer(params...)
	return errors.New(replacer.Replace(format))
}
