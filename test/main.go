package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Zhan9Yunhua/blog-svr/utils"
)

func main() {
	// ch := make(chan int, 1)
	type A struct {
		Uid string
		A   string
		b   int
		d   bool
		C   map[string]interface{}
	}

	type B struct {
		Uid string
		A   string
		b   int
		d   bool
		C   map[string]interface{}
	}

	a := &A{Uid: "aaa", A: "bbb", C: map[string]interface{}{"test": 1}}
	b := &B{}
	fmt.Println(utils.StructCopy(a, b))
	// StructCopy(a, b)
	fmt.Println(b)
}
func StructCopy(src, dst interface{}) error {
	srcV, err := srcFilter(src)
	if err != nil {
		return err
	}
	dstV, err := dstFilter(dst)
	if err != nil {
		return err
	}
	srcKeys := make(map[string]bool)
	for i := 0; i < srcV.NumField(); i++ {
		srcKeys[srcV.Type().Field(i).Name] = true
	}
	for i := 0; i < dstV.Elem().NumField(); i++ {
		fName := dstV.Elem().Type().Field(i).Name
		if _, ok := srcKeys[fName]; ok {
			v := srcV.FieldByName(dstV.Elem().Type().Field(i).Name)
			if v.CanInterface() {
				dstV.Elem().Field(i).Set(v)
			}
		}
	}
	return nil
}
func srcFilter(src interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(src)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Zero(v.Type()), errors.New("src type error: not a struct or a pointer to struct")
	}
	return v, nil
}
func dstFilter(src interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(src)
	if v.Type().Kind() != reflect.Ptr {
		return reflect.Zero(v.Type()), errors.New("src type error: not a pointer to struct")
	}
	if v.Elem().Kind() != reflect.Struct {
		return reflect.Zero(v.Type()), errors.New("src type error: not point to struct")
	}
	return v, nil
}
