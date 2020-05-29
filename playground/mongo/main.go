package main

import (
	"log"
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
	if err := database.Create(&info); err != nil {
		panic(err)
	}
	info.Name = "Update"
	if err := database.Update(&info); err != nil {
		panic(err)
	}
	rInfo := new(Info)
	if err := database.Get(&rInfo, db.Query{"Id", "=", 3}); err == nil {
		if rInfo.Name != info.Name {
			panic("no equal")
		}
		log.Println(rInfo)
	}
}
