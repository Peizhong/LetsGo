package main

import (
	"log"
	"reflect"
	"unsafe"
)

type people struct {
	Name string
	Age  int
}

func WhatPointer() {
	me := people{Name: "wpz", Age: 29}
	pName := (*string)(unsafe.Pointer(&me))
	*pName = "wulala"
	offsetAge := uintptr(unsafe.Pointer(&me)) + unsafe.Offsetof(me.Age)
	pAge := (*int)(unsafe.Pointer(offsetAge))
	*pAge = 20
	log.Print(me)

	var a, b struct{}
	println(a == b, &a, &b, &a == &b, unsafe.Pointer(&a) == unsafe.Pointer(&b), uintptr(unsafe.Pointer(&a)) == uintptr(unsafe.Pointer(&b)))

	var c, d people
	println(c == d, &c, &d, &c == &d, unsafe.Pointer(&c) == unsafe.Pointer(&d), uintptr(unsafe.Pointer(&c)) == uintptr(unsafe.Pointer(&d)))
}

func dump(v interface{}) {
	raw := reflect.ValueOf(v)
	switch raw.Kind() {
	case reflect.Struct:
		println("struct")
	case reflect.Ptr:
		println("prt")
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		println("array or slice")
	}
	// 结构体
	// 指针
	// 结构体数组/切片
	// 指针数组/切片
}

func WhatReflect() {
	me := people{Name: "wpz", Age: 29}
	v := reflect.ValueOf(me)
	name := v.FieldByName("Name")
	println(v.String(), name.String())

	pv := reflect.ValueOf(&me)
	v = reflect.Indirect(pv)
	name = v.FieldByName("Name")
	println(v.String(), name.String())

	a := [1]int{1}
	s := []int{}
	dump(a)
	dump(s)
}
