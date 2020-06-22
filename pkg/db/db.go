package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peizhong/letsgo/pkg/config"
)

func raw() {
	db, err := sql.Open("", "")
	if err != nil {
		return
	}
	db.Ping()

	r,err:=db.Query("", "")
	r.Close()
}

func getDBConnString(db string) (dbType, connStr string) {
	dbType = db
	switch db {
	case "mysql":
		connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	case "sqlite3":
		connStr = "letsgo.db"
	case "mongo":
		connStr = fmt.Sprintf("mongodb://%s:27017", config.MongoDBHost)
	default:
		panic("no database")
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
		return &GormHandler{dBType: dbType, context: context.Background()}
	case "mongo":
		return &MongoHandler{}
	}
	return nil
}

type Query struct {
	Key   string
	Op    string
	Value interface{}
}

type ORMHandler interface {
	SetContext(string, interface{})
	Ping() error
	Create(interface{}) error
	Get(interface{}, ...Query) error
	Gets(interface{}, ...Query) (int, error)
	Update(interface{}) error
}
