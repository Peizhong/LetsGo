package main

import "reflect"

type i1 interface {
	Do1()
}

type i2 interface {
	i1
	Do2()
}

type One struct {
}

func (o *One) Do1() {

}

type Two struct {
	One
}

func (t *Two) Do1() {
	var i i2 = t
	i.Do1()
}

func (t *Two) Do2() {
	var i i2 = t
	i.Do1()
}

func DoInterface() {

	tv := &Two{}
	v := reflect.ValueOf(tv)
	rv := reflect.Indirect(v)
	ev := v.Elem()
	if rv == ev {

	}
	v.MethodByName("a").Call([]reflect.Value{reflect.ValueOf(1)})

	c1 := make(chan int)
	c2 := make(chan int)
	if c1 != c2 {

	}
	s1 := []int{}
	/// s2 := []int{}
	if reflect.TypeOf(s1).Comparable() {
		// false
	}
	m1 := make(map[int]int)
	// m2 := make(map[int]int)
	if reflect.TypeOf(m1).Comparable() {
		// false
	}
	reflect.DeepEqual(s1, m1)
}
