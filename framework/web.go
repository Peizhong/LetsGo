package framework

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type DbContext struct {
	//
	GetDatabase func() (*gorm.DB, error)
	GetCache    func() (redis.Conn, error)
}
