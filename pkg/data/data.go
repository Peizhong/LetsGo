package data

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"reflect"
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
	context.Background()
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

func GetTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func GetMap(i interface{}) (string,map[string]interface{}) {
	table := GetTypeName(i)
	m := map[string]interface{}{}
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		m[t.Field(i).Name] = f.Interface()
	}
	return table, m
}

func GetMapAsJson(i interface{}) (string,map[string]interface{}) {
	table := GetTypeName(i)
	m := map[string]interface{}{}
	data, _ := json.Marshal(i)
	json.Unmarshal(data, &m)
	return table, m
}
