package main

import (
	"time"

	"github.com/peizhong/letsgo/pkg/db"
)

type Info struct {
	Id    int `pk:"true"`
	Value string
	Name  string
	Time  time.Time
}

func main() {
	database := db.DBFactory("mongo")
	if err := database.Ping(); err != nil {
		panic(err)
	}
	info := Info{3, "300", "Time", time.Now()}
	if err := database.Update(&info); err != nil {
		panic(err)
	}
}
