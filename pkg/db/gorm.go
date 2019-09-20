package db

import "github.com/jinzhu/gorm"

type GormHandler struct {
	dBType, connStr string
}

func (g *GormHandler) Do(f func(*gorm.DB) (interface{}, error)) (interface{}, error) {
	g.dBType, g.connStr = GetDBConnString("")
	if db, err := gorm.Open(g.dBType, g.connStr); err != nil {
		return nil, err
	} else {
		db.LogMode(true)
		r, err := f(db)
		return r, err
	}
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

func (g *GormHandler) Get(i interface{}) error {
	_, err := g.Do(func(db *gorm.DB) (interface{}, error) {
		db.First(i)
		return nil, nil
	})
	return err
}

func (g *GormHandler) Gets(i interface{}) error {
	_, err := g.Do(func(db *gorm.DB) (interface{}, error) {
		db.Find(i)
		return nil, nil
	})
	return err
}