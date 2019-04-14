package main

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

/*
consul for service discovery
*/

func init() {

}

func main() {
	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err.Error())
	}
	defer db.Close()
}
