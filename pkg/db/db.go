package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/peizhong/letsgo/pkg/config"
)

func GetConnString() (dbType, connStr string) {
	dbType = config.DBType
	switch dbType {
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	case "sqlite3":
		connStr = "letsgo.db"
	}
	return
}

func Once(f func(*gorm.DB) (interface{}, error)) (interface{}, error) {
	dbType, connStr := GetConnString()
	if db, err := gorm.Open(dbType, connStr); err != nil {
		return nil, err
	} else {
		db.LogMode(true)
		r, err := f(db)
		return r, err
	}
}

func Migrate(i interface{}) error {
	_, err := Once(func(db *gorm.DB) (interface{}, error) {
		db.AutoMigrate(i)
		return nil, nil
	})
	return err
}

func Create(i interface{}) error {
	_, err := Once(func(db *gorm.DB) (interface{}, error) {
		db.Create(i)
		return nil, nil
	})
	return err
}

func Get(i interface{}) error {
	_, err := Once(func(db *gorm.DB) (interface{}, error) {
		db.First(i)
		return nil, nil
	})
	return err
}

func Gets(i interface{}) error {
	_, err := Once(func(db *gorm.DB) (interface{}, error) {
		db.Find(i)
		return nil, nil
	})
	return err
}
