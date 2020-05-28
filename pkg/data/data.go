package data

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
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

// get tag from single field
func GetTag(i interface{}, tag string) (string, string) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", ""
	}
	for n := 0; n < t.NumField(); n++ {
		field := t.Field(n)
		if value, ok := field.Tag.Lookup(tag); ok {
			return field.Name, value
		}
	}
	return "", ""
}

func GetMap(i interface{}) (string, map[string]interface{}) {
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

func GetMapAsJson(i interface{}) (string, map[string]interface{}) {
	table := GetTypeName(i)
	m := map[string]interface{}{}
	data, _ := json.Marshal(i)
	json.Unmarshal(data, &m)
	return table, m
}

// primarykey is tag:pk, i should be pointer
func GetPrimaryKey(i interface{}) map[string]interface{} {
	if filed, isKey := GetTag(i, "pk"); isKey == "true" {
		value := reflect.ValueOf(i).Elem().FieldByName(filed).Interface()
		return map[string]interface{}{filed: value}
	}
	return nil
}
