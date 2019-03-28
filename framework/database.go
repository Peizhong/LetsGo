package framework

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	connStr, _ := GetConnectionString("avmt")
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
	fmt.Printf("The classify name of 10012 is %v", classifyName)
}
