package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Str interface {
	String() string
}

type T struct {
	TName string `validator:"number"`
	TAge  int    `validator:"number"`
}

func (t *T) String() string {
	fmt.Println("T")
	return t.TName
}

type P struct {
	Name string            `validator:"required"`
	Age  int               `validator:"number=[11.1|20]"`
	Sex  map[string]string `validator:"multi={name:zz,age:1}"`
	Str
}

func (p *P) String() string {
	fmt.Println("P")
	fmt.Println(p.Str.String())
	return p.Name
}

type Person struct {
	Name string   `validator:"required||string=[2|_" map:"name"` // 必填。2<=len
	Age  int      `validator:"number=10|20]" map:"age"`          // 选填。10<Age<=20 / 0
	Sex  int      `validator:"number||in=0,1,2"`                 // 选填。值只能是0||1||2
	Car  []string `validator:"multi=_|5]||in=LEXUS,BMW"`         // 选填。len>=5且包含LEXUS||BMW
}

func main() {
	p := new(P)
	p.Str = &T{}

	fmt.Println(p.String())
	// p := Person{Name: "z", Age: 11, Sex: 1, Car: []string{"AUDI"}}
	//
	// vali := validator.NewValidator()
	// err := vali.Validate(p)
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
