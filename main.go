package main

import (
	"fmt"
	// . fmt 省略前缀 Println("hello world")
	// f fmt 重命名 f.Println("hello world")
	// _ "github.com/ziutek/mymysql/godrv" 引入该包, init()，而不直接使用包里面的函数

	"log"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/peizhong/letsgo/framework"
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
	//customName = defaultName        //not allowed

	_, _ = defaultName, customName
}

func doRelect() {
	joe := framework.Person{
		Name:    "wang peizhong",
		Address: "shezhen",
	}
	pJoe := &joe
	pJoe.WhoAmI("GH")
	checkReflect(pJoe)
}

func checkReflect(obj interface{}) {
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

// <-chan writeonly, chan<- readonly
func helloProducer(ch chan<- int) {
	fmt.Println("producer started")
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("write to channel", i)
		ch <- i
	}
	close(ch)
	fmt.Println("producer all done")
}

// cost of creating a Goroutine is tiny when compared to a thread
func doGoroutines(count int32) {
	ch := make(chan int)
	// whern a goroutine is started, the goroutine call returns immediately to the next line of code after the Goroutine call
	go helloProducer(ch)
	// write and send to channel are bolcking by default
	// v, ok := <-ch
	for v := range ch {
		fmt.Println("Received ", v)
	}
}

// 如果导入了多个包，先初始化包的参数，然后init()，最后执行package的main()
func init() {
	// 将日志输出到标准输出
	log.SetOutput(os.Stdout)
}

// 每个package必须有个main
func main() {
	doGoroutines(123)

	doRelect()

	//play.FromSQL2NoSQL("‪C:/Users/wxyz/Desktop/avmt.db", "", "")
	//search.Run("mimi")

	framework.RunServer(8080)
}
