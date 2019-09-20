package db

import (
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
	n1 := GetTableName(s)
	n2 := GetTableName(s[0])
	n3 := GetTableName(sp)
	n4 := GetTableName(sp[0])
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