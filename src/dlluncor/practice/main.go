package main

import (
	"fmt"
	"reflect"
)

func MakeMap(k0 reflect.Type, v0 reflect.Type) reflect.Value {
	b := reflect.MakeMap(reflect.MapOf(k0, v0))
	return b
}

func val(in interface{}) reflect.Value {
	return reflect.ValueOf(in)
}

func Ex0() {
	// Creating generic maps essentially.
	a := map[string]int{
		"hi": 2,
	}
	b := MakeMap(reflect.TypeOf(""), reflect.TypeOf(1))
	// Add to the map some values.
	b.SetMapIndex(val("hi"), val(2))
	bb := b.Interface()
	d := bb.(map[string]int)
	c := reflect.DeepEqual(a, d)
	fmt.Printf("%v, %v, %v\n", reflect.TypeOf(a), reflect.TypeOf(b), d)
	fmt.Printf("%v\n", c)
}

func main() {
	fmt.Printf("Practice.\n")
	Ex0()
}
