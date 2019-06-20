package data

import (
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"strconv"
)

func IntTryParse(s string) (n int, b bool) {
	if num, err := strconv.Atoi(s); err == nil {
		return num, true
	}
	return
}

func Int64TryParse(s string) (n int64, b bool) {
	if num, err := strconv.ParseInt(s, 10, 64); err == nil {
		return num, true
	}
	return
}

func NewGuid() string {
	id := uuid.New().String()
	return id
}

func GetJsonValue(data, path string) string {
	value := gjson.Get(data, path)
	return value.String()
}
