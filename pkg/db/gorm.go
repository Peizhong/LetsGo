package db

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	tracinglog "github.com/opentracing/opentracing-go/log"
)

type GormHandler struct {
	dBType, connStr string
	db              *gorm.DB
	context         context.Context
}

const (
	tracing_span_name = "gorm_sql"
	tracing_span      = "sql_span"
)

var onceMysql sync.Once

func (g *GormHandler) SetContext(key string, val interface{}) {
	// tracing 污染代码
	if key == "" {
		if ctx, ok := val.(context.Context); ok {
			g.context = ctx
		}
	} else {
		g.context = context.WithValue(g.context, key, val)
	}
}

func (g *GormHandler) Do(f func(*gorm.DB) (interface{}, error)) (interface{}, error) {
	onceMysql.Do(func() {
		var err error
		g.dBType, g.connStr = getDBConnString(g.dBType)
		g.db, err = gorm.Open(g.dBType, g.connStr)
		if err != nil {
			panic(err)
		}
		g.db.LogMode(false)
		g.db.Callback().Query().Before("gorm:query").Register("tracing:query_before", g.before)
		g.db.Callback().Query().After("gorm:after_query").Register("tracing:query_after", g.after)
		g.db.DB().SetMaxIdleConns(10)
		g.db.DB().SetMaxOpenConns(100)
		// g.db.Callback().Create()
		// g.db.Callback().Update()
		// g.db.Callback().RowQuery()
	})
	r, err := f(g.db)
	return r, err
}

func (g *GormHandler) before(scope *gorm.Scope) {
	// context设置了
	span, _ := opentracing.StartSpanFromContext(g.context, tracing_span)
	scope.DB().Set(tracing_span, span)
}

func (g *GormHandler) after(scope *gorm.Scope) {
	if sp, ok := scope.Get(tracing_span); ok {
		span := sp.(opentracing.Span)
		defer span.Finish()
		span.LogFields(tracinglog.String("sql", scope.SQL))
		if scope.HasError() {
			span.LogFields(tracinglog.Error(scope.DB().Error))
		}
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

func (g *GormHandler) Update(i interface{}) error {
	return nil
}
