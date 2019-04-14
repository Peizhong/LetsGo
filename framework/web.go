package framework

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type DbContext struct {
	Database *gorm.DB
	Cache    redis.Conn
}
