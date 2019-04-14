package framework

import (
	"github.com/jinzhu/gorm"
)

func NewMySQLConn(connString string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return db, nil
}
