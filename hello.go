package main

import (
	"fmt"
	// . fmt 省略前缀 Println("hello world")
	// f fmt 重命名 f.Println("hello world")
	// _ "github.com/ziutek/mymysql/godrv" 引入该包, init()，而不直接使用包里面的函数
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func test() {
	// 定义变量
	var b bool
	// 定义并赋值
	var n = 1
	// 简短声明
	_, v1, v2, v3 := 0, 1, 3, 5
	// 常量
	const value0 = 1234
	value1 := 1
	_, value2 := 2, value1
	// 内置数据类型
	var enable = true
	// 浮点
	num := 1.1
	var helloString = "how you doing"
	code := http.StatusOK
	// 字符串转换为[]byte
	c := []byte(helloString)
	s2 := string(c)
	// 字符串切片
	s3 := "h" + s2[1:]
	// 没有支付转义
	sraw := s3 + `ads
	ads
	`
	// array
	nums := [4]int{1, 3, 5, 7}
	// 省略长度
	cnums := [...]int{2, 4, 6, 8}
	// 二维数组 row*column
	doubleArray := [...][4]int{{1, 2, 3, 4}, {5, 6, 7, 8}}
	// slice 不固定长度，引用类型
	slice := []byte{'a', 'b', 'c', 'd'}

	var array [10]int
	// 长度2，容量2-10:8
	slice2 := array[2:4]
	// append函数会改变slice所引用的数组的内容，从而影响到引用同一数组的其它slice。
	// 但当slice中没有剩余空间（即(cap-len) == 0）时，此时将动态分配新的数组空间。返回的slice数组指针将指向这个空间，而原数组的内容将保持不变；其它引用此数组的slice则不受影响
	slice3 := append(slice2, 1)
	// iota枚举， const中重置，每行加1
	const (
		red   = iota //0
		green = iota //1
		blue  = iota //2
	)
	// 长度2，容量6
	slice4 := make([]byte, 2, 6)
	// map, 用make初始化
	numbers := map[string]string{}
	numbers["one"] = "a"
	numbers["two"] = "b"
	//value2, enable := numbers["one2"] 生成new value
	helloString, enable = numbers["one2"]

	// make用于内建类型（map、slice 和channel）的内存分配。new用于各种类型的内存分配。

	fmt.Println("%s", sraw)
	fmt.Println(code, b, n, v1, v2, v3, value2, enable, num, nums, cnums, doubleArray, slice, slice3, slice4)
	fmt.Println("!oG ,olleH")
	// 大写变量/函数是可导出的，其他包可以读取，小写的是私有

	logic()
	readWrite()
	callFunc(100, realFunc)
}

func logic() {
	sum := 0
	for index := 0; index < 10; index++ {
		sum += index
	}
	// 相当于while了
	for sum > 0 {
		sum -= 1
	}
	mp := map[int]string{}
	mp[1] = "1"
	mp[2] = "2"
	kv := ""
	for k, v := range mp {
		sum += k
		kv += v
	}
	i := 10
	switch i {
	case 1:
		sum = 1
	case 2:
		sum = 2
	default:
		sum = 0
	}
	aa := 123
	bb := addPtr(&aa)
	fmt.Println(sum, bb)
}

func multiReturn(a int, b int) (sa string, sb string) {
	sa = "123"
	sb = "321"
	return
}

func multiParams(arg ...int) (res int) {
	// arg 是个 slice
	for _, n := range arg {
		res += n
	}
	return res
}

// 传指针
func addPtr(a *int) int {
	*a = *a + 1
	return *a
}

// 延迟
func readWrite() bool {
	fp, err := os.Open("README.md")
	if err != nil {
		fmt.Println("Open file error: ", err)
		return false
	}
	// defer后指定的函数会在函数退出前调用，后进先出模式
	defer fp.Close()
	return true
}

// 声明了一个函数类型
type testFunc func(int) bool

func realFunc(a int) bool {
	if z := a % 2; z == 1 {
		return true
	}
	return false
}

func callFunc(a int, f testFunc) (r bool) {
	r = f(a)
	return
}

type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02"`
}

func userInfo(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(person.Name)
	log.Println(person.Address)
	log.Println(person.Birthday)

	c.String(200, "Success")
}

// simulate some private data
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

// 如果导入了多个包，先初始化包的参数，然后init()，最后执行package的main()
func init() {

}

// 每个package必须有个main
func main() {
	test()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets endpoint
	// hit "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST("/user", userInfo)
	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	// Multiple files
	router.POST("/uploads", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	// cookie
	router.GET("/cookie", func(c *gin.Context) {

		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		fmt.Printf("Cookie value: %s \n", cookie)
	})
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(":8080")
}
