package db

import (
	"github.com/stretchr/testify/assert"
	"letsgo/models/catalog"
	"testing"
)

func TestMigrate(t *testing.T) {
	classify := &catalog.Classify{}
	err := Migrate(classify)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	classify := &catalog.Classify{}
	err := Get(classify)
	assert.Nil(t, err)
}

func TestCreate(t *testing.T) {
	classify := &catalog.Classify{
		Name: "wulala",
	}
	err := Create(classify)
	assert.Nil(t, err)
}
