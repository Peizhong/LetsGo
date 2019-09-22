package db

import (
	"github.com/peizhong/letsgo/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

type demo struct {
	K string
	V string
}

func TestGetTableName(t *testing.T) {
	s := []demo{
		{"1", "1"},
	}
	sp := []*demo{
		{"1", "1"},
	}
	n1 := data.GetTypeName(s)
	n2 := data.GetTypeName(s[0])
	n3 := data.GetTypeName(sp)
	n4 := data.GetTypeName(sp[0])
	assert.Equal(t, n1, n2, n3, n4)
}

func TestMongoHandler_Create(t *testing.T) {
	storage := DBFactory("mongo")
	err := storage.Ping()
	assert.NoError(t, err)
	s := []demo{
		{"1", "1"},
	}
	err = storage.Create(s[0])
	assert.NoError(t, err)
}

func TestMongoHandler_Get(t *testing.T) {
	s := demo{K: "1", V: "1"}
	q := demo{V: "1"}
	storage := DBFactory("mongo")
	err := storage.Get(&q)
	assert.NoError(t, err)
	assert.Equal(t, s.V, q.V)
}

func TestMongoHandler_Gets(t *testing.T) {
	r := []demo{}
	storage := DBFactory("mongo")
	_, err := storage.Gets(&r)
	assert.NoError(t, err)
}
