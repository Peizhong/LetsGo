package db

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
)

type GormHandler struct {
	dBType, connStr string
	db              *gorm.DB
}

var onceMysql sync.Once

func (g *GormHandler) Do(f func(*gorm.DB) (interface{}, error)) (interface{}, error) {
	onceMysql.Do(func() {
		var err error
		g.dBType, g.connStr = getDBConnString(g.dBType)
		g.db, err = gorm.Open(g.dBType, g.connStr)
		if err != nil {
			panic(err)
		}
		g.db.LogMode(false)
	})
	r, err := f(g.db)
	return r, err
}

func (g *GormHandler) Ping() error {
	_, err := g.Do(func(db *gorm.DB) (interface{}, error) {
		return nil, db.DB().Ping()
	})
	return err
}

func (g *GormHandler) Migrate(i interface{}) error {
	_, err := g.Do(func(db *gorm.DB) (interface{}, error) {
		db.AutoMigrate(i)
		return nil, nil
	})
	return err
}

func (g *GormHandler) Create(i interface{}) error {
	_, err := g.Do(func(db *gorm.DB) (interface{}, error) {
		db.Create(i)
		return nil, nil
	})
	return err
}

func (*GormHandler) buildQuery(query ...Query) (querySql string, args []interface{}) {
	var qs []string
	for _, q := range query {
		qs = append(qs, fmt.Sprintf("%s %s ?", q.Key, q.Op))
		args = append(args, q.Value)
	}
	querySql = strings.Join(qs, " AND ")
	return
}

func (g *GormHandler) Get(i interface{}, q ...Query) error {
	_, err := g.Do(func(db *gorm.DB) (interface{}, error) {
		t := reflect.TypeOf(i).Elem()
		db.Table(t.Name()).Where(g.buildQuery(q...)).First(i)
		return nil, nil
	})
	return err
}

func (g *GormHandler) Gets(i interface{}, q ...Query) (int, error) {
	_, err := g.Do(func(db *gorm.DB) (interface{}, error) {
		db.Find(i)
		return nil, nil
	})
	return 0, err
}
