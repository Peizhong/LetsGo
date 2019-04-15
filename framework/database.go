package framework

import (
	"github.com/jinzhu/gorm"
)

func NewMySQLConn(connString string) func() (*gorm.DB, error) {
	return func() (*gorm.DB, error) {
		// does not establish any connections to the database
		db, err := gorm.Open("mysql", connString)
		if err != nil {
			return nil, err
		}
		db.LogMode(true)
		return db, nil
	}
}
