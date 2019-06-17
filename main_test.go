package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
	"unsafe"

	"letsgo/framework/log"
)

// 闭包
func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func doVariables() {
	var width, height int
	var (
		name = "naveen"
		age  = 29
	)
	// least one of the variables in the left side of := is newly declared
	sex, width, height := "male", 174, 65
	_, _, _, _, _ = name, age, sex, width, height
}

func doTypes() {
	// bool
	a := true
	b := false
	c := a && b
	d := a || b
	e := unsafe.Sizeof(a)
	_, _, _ = c, d, e

	// float and complex
	f := 5.67
	g := complex(5, 7)
	g = 8 + 27i
	_, _ = f, g

	// byte = unint8, rune = int32

	// type conversion: no automatic type promotion or conversion
	i := 55
	j := 6.78
	// k := i + j // int + float not allowed
	sum0 := i + int(j)
	sum1 := float64(i) + j
	_, _ = sum0, sum1
}

func doConstants() {
	const a = 55
	// const b = math.Sqrt(3) // value of a constant should be known at compile time

	// string constant like "hello world" does not have any type
	// untyped constants have a default type associated with them and they supply it if and only if a line of code demands it
	const hello = "hello world"

	var defaultName = "Sam" //allowed
	type myString string
	var customName myString = "Sam" //allowed
	//customName = defaultName        //different type, not allowed
	fmt.Println(fmt.Sprintf("%T %T", defaultName, customName))
}

func TestCondition(t *testing.T) {
	for i := 0; i < 10; i++ {
		if i < 5 {
			continue
		}
		if i == 6 {
			break
		}
	}
	j := 0
	// aka while
	for j < 10 {
		j++
	}
	letter := "i"
	switch letter {
	case "a", "e", "i", "o", "u": //multiple expressions in case
		fmt.Println("vowel")
	default:
		fmt.Println("not a vowel")
	}
	// Expressionless switch
	// If the expression is omitted, the switch is considered to be switch true and each of the case expression is evaluated for truth
	num := 50
	num += 2
	switch {
	case num >= 0 && num <= 50:
		num += 2
		fmt.Println("num is greater than 0 and less than 50")
	case num >= 51 && num <= 100:
		fmt.Println("num is greater than 51 and less than 100")
		fallthrough //go to next
	case num >= 101:
		fmt.Println("num is greater than 100")
	}
}

func doSlice() {
	a := [...]int{1, 2, 3, 4}
	s := a[:3]
	// 如果slice还有容量，会把原来array的值改变
	s = append(s, 5)

	cars := []string{"Ferrari", "Honda", "Ford"}
	// The capacity of the new slice is twice that of the old slice
	cars = append(cars, "Toyota")

	countries := []string{"USA", "Singapore", "Germany", "India", "Australia"}
	neededCountries := countries[:len(countries)-2]
	countriesCpy := make([]string, len(neededCountries))
	copy(countriesCpy, neededCountries) //copies neededCountries to countriesCpy, min len

}

func doMap() {
	personSalary := map[string]int{
		"steve": 12000,
		"jamie": 15000,
	}
	personSalary["mike"] = 9000
	if _, ok := personSalary["joe"]; !ok {
		personSalary["joe"] = 14000
	}
	// order of the retrieval of values from a map when using for range is not guaranteed
	for k, v := range personSalary {
		println(k, "", v)
	}
	delete(personSalary, "joe")
}

func changeString(s []rune) string {
	s[0] = 'a'
	return string(s)
}

func doString() {
	name := "Hello World"
	runes := []rune(name)
	str := string(runes)
	fmt.Printf("length of %s is %d\n", str, utf8.RuneCountInString(str))

	name = changeString([]rune(name))
}

func doStructure() {
	// anonymous structures
	emp3 := struct {
		firstName, lastName string
		age, salary         int
	}{
		firstName: "Andreah",
		lastName:  "Nikola",
		age:       31,
		salary:    5000,
	}
	emp4 := struct {
		firstName, lastName string
		age, salary         int
	}{
		firstName: "Andreah",
		lastName:  "Nikola",
		age:       31,
		salary:    5000,
	}
	if emp3 == emp4 {

	}
	// Structs are value types and are comparable if each of their fields are comparable
	// Struct variables are not comparable if they contain fields which are not comparable, like map
}

func doMethod() {
	/*
		pointer receiver and when to use value receiver
		Pointers receivers can also be used in places where it's expensive to copy a data structure.
		Consider a struct which has many fields.
		Using this struct as a value receiver in a method will need the entire struct to be copied which will be expensive.
		In this case if a pointer receiver is used, the struct will not be copied and only a pointer to it will be used in the method.
		In all other situations value receivers can be used.
	*/
}

func checkReflect(obj interface{}) {
	// compare the concrete type of an interface
	switch obj.(type) {
	case string:
	case int:
	default:
	}
	getType := reflect.TypeOf(obj)
	getKind := getType.Kind()
	getValue := reflect.ValueOf(obj)
	fmt.Println(getType, getKind, getValue)
	switch getKind {
	case reflect.Struct:
		for i := 0; i < getType.NumField(); i++ {
			field := getType.Field(i)
			value := getValue.Field(i).Interface()
			fmt.Printf("Field:%d %s: %v = %v\n", i, field.Name, field.Type, value)
		}
		// NumMethod公共方法
		for i := 0; i < getType.NumMethod(); i++ {
			m := getType.Method(i)
			fmt.Printf("%s: %v\n", m.Name, m.Type)
		}
	case reflect.Ptr:
		// 传入的指针，根据elem获得指向的值
		realValue := getValue.Elem()
		realType := realValue.Type()
		for i := 0; i < realType.NumField(); i++ {
			field := realType.Field(i)
			value := realValue.Field(i).Interface()
			fmt.Printf("Field:%d %s: %v = %v\n", i, field.Name, field.Type, value)
		}
		for i := 0; i < realType.NumMethod(); i++ {
			m := realType.Method(i)
			fmt.Printf("%s: %v\n", m.Name, m.Type)
		}
		if _, exist := realType.MethodByName("WhoAmI"); exist {
			methodValue := realValue.MethodByName("WhoAmI")
			args := []reflect.Value{reflect.ValueOf("GH")}
			methodValue.Call(args)
		}
	default:

	}
}

type job struct {
	// unsafe 对齐
	b        byte
	id       int
	randomno int
	i32      int32
	i64      int64
}

func createBeater(noOfJobs int, jobs chan<- job) {
	for i := 0; i < noOfJobs; i++ {
		randomno := rand.Intn(999)
		j := job{1, i, randomno, 1, 2}
		fmt.Println("go ", i)
		jobs <- j
	}
	close(jobs)
}

func worker(wg *sync.WaitGroup, jobs <-chan job) {
	for range jobs {
		time.Sleep(1 * time.Second)
		<-jobs
	}
	wg.Done()
}

func createWorker(noOfWorkers int, jobs chan job) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg, jobs)
	}
	wg.Wait()
}

// worker pool is a collection of threads which are waiting for tasks to be assigned to them.
// Once they finish the task assigned, they make themselves available again for the next task
func doWorkerPool() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	jobs := make(chan job, cpus)
	go createBeater(100, jobs)
	createWorker(cpus, jobs)
}

type C1 struct {
	context.Context
}

func (c C1) Run() {
	log.Info("waiting for parent to stop it")
	<-c.Done()
	log.Info("child completed")
}

func TestContext(t *testing.T) {
	pt := context.Background()
	ct1 := context.WithValue(pt, "hello", "world")
	ct11 := context.WithValue(ct1, "world", "hello")
	v := ct11.Value("hello").(string)
	_ = v
	ct2, cancel := context.WithCancel(ct1)
	ct0 := C1{
		Context: ct2,
	}
	go func() {
		ct0.Run()
	}()
	time.AfterFunc(time.Second*3, func() {
		cancel()
	})
	time.Sleep(time.Second * 5)
}

func f2i(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

type company struct {
	j    job
	name int32
}

func TestPtr(t *testing.T) {
	var f float64
	f = 12
	log.Info(unsafe.Sizeof(f))
	var d uint64
	log.Info(unsafe.Sizeof(d))
	d = f2i(f)
	j := &job{
		id:       1,
		randomno: 2,
	}
	ao := unsafe.Alignof(j.b)
	os := unsafe.Offsetof(j.b)
	so := unsafe.Sizeof(j.b)
	// 结构体对齐值，
	ao = unsafe.Alignof(*j)
	// 64位对齐系数=8
	so = unsafe.Sizeof(*j)
	cp := &company{
		*j,
		123,
	}
	ao = unsafe.Alignof(*cp)
	so = unsafe.Sizeof(*cp)
	os = unsafe.Offsetof(cp.name)
	log.Info(ao, os, so)
	str := "aaaa"
	so = unsafe.Sizeof(str)
	so = unsafe.Sizeof(&str)
	// array, slice 的区别
	a := [...]int{1, 2, 3}
	prtA := unsafe.Pointer(&a)
	prt0 := unsafe.Pointer(&a[0])
	sum := uintptr(prt0) - uintptr(prtA)
	log.Info(sum)
	/*
	   （1）任何类型的指针都可以被转化为Pointer
	   （2）Pointer可以被转化为任何类型的指针
	   （3）uintptr可以被转化为Pointer
	   （4）Pointer可以被转化为uintptr
	*/
	var x struct {
		a bool
		b int16
		c []int
	}
	fmt.Println(x.b) // "42"
	tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	// gc可能移动了x，导致tmp uintptr指向的地址无效
	pb := (*int16)(unsafe.Pointer(tmp))
	*pb = 42
	var p1 struct {
		a bool
		b int32
		c int8
		d int64
		e byte
	}
	var p2 struct {
		a bool  //0 1
		e byte  //1 1
		c int8  //2 1 + 1
		b int32 //4 4
		d int64 //8 8
	}
	var p3 struct {
		a bool
	}
	fmt.Println(unsafe.Alignof(p3)) //64位不超过8
	assert.Equal(t, unsafe.Sizeof(p1), uintptr(32))
	assert.Equal(t, unsafe.Sizeof(p2), uintptr(16))

	fmt.Println(x.b) // "42"
	log.Info(unsafe.Sizeof(j))
	log.Info(unsafe.Alignof(j.b))
	log.Info(unsafe.Sizeof(*j))
	log.Info(unsafe.Sizeof(j.i32))
	log.Info(unsafe.Sizeof(j.i64))
}

func slowCal(fac int) int {
	if fac < 2 {
		return fac
	}
	return slowCal(fac-1) + slowCal(fac-2)
}
