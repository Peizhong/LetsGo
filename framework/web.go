package framework

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type GoContext struct {
	Context     context.Context
	GetDatabase func() (*gorm.DB, error)
	GetCache    func() (redis.Conn, error)
}
