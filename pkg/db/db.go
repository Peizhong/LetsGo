package db

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/peizhong/letsgo/pkg/config"
)

func GetDBConnString(db string) (dbType, connStr string) {
	dbType = db
	if dbType == "" {
		dbType = config.DBType
	}
	switch dbType {
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	case "sqlite3":
		connStr = "letsgo.db"
	case "mongo":
		connStr = fmt.Sprintf("mongodb://%s:27017", config.DBHost)
	}
	return
}

func DBFactory(db string) ORMHandler {
	dbType := db
	if dbType == "" {
		dbType = config.DBType
	}
	switch dbType {
	case "mysql", "sqlite3":
		return &GormHandler{}
	case "mongo":
		return &MongoHandler{}
	}
	return nil
}

type ORMHandler interface {
	Ping() error
	Create(interface{}) error
}
