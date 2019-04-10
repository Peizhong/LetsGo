package framework

import (
	"database/sql"
	"fmt"

	// linq
	linq "github.com/ahmetb/go-linq"
	"github.com/jinzhu/gorm"
)

type dbEntity interface{}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// If a struct type starts with a capital letter, then it is a exported type and it can be accessed from other packages
// Similarly if the fields of a structure start with caps, they can be accessed from other packages
type dbContext struct {
	connectionString string
	dbSet            []dbEntity
}

func initDatabase() {
	connStr := GetConnectionString("avmt")
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	stmtOut, err := db.Prepare("SELECT CLASSIFY_NAME FROM dm_classify where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()
	var classifyName string
	err = stmtOut.QueryRow(10012).Scan(&classifyName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(fmt.Sprintf("The classify name of 10012 is %v", classifyName))

	demo := dbContext{
		dbSet: []dbEntity{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
	}
	var q []int
	linq.From(demo.dbSet).Where(func(c interface{}) bool {
		// Type Assertion: get the underlying value of interface
		// v, ok := i.(T), if i is not T then ok will be false and v will have the zero value of type T, no panic
		return c.(int) > 10
	}).Select(func(c interface{}) interface{} {
		return c.(int)
	}).ToSlice(&q)
}

func SayHiOrm() {
	connStr := GetConnectionString("avmt")
	// 已在main中注册了mysql的驱动，自带连接池?
	db, err := gorm.Open("mysql", connStr)
	db.LogMode(true)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	// db.AutoMigrate(&Product{})

	// Read
	var product = &Product{}
	db.First(product, 12)                  // find product with id 1
	db.First(product, "code = ?", "L1212") // find product with code l1212

	if product.ID == 0 {
		// Create
		db.Create(&Product{Code: "L1212", Price: 1000})
	} else {
		// Update - update product's price to 2000
		db.Model(&product).Update("Price", 2000)
		// Delete - delete product
		db.Delete(&product)
	}
}
