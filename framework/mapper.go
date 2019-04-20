package framework

import (
	"encoding/json"
	"log"
	"reflect"
	"time"
)

type ObjOne struct {
	Key   string
	Value string
}

type ObjTwo struct {
	Key   string
	Value string
}

type DemoOne struct {
	hidden     int
	Id         int
	Value      string
	UpdateTime time.Time
	Data       []*ObjOne
	Nums       []int
}

type DemoTwo struct {
	Id         int
	Value      string
	UpdateTime time.Time
	Data       []*ObjTwo
	Nums       []int
}

func MapTo(src, dst interface{}) {
	// UseDefault
	// IngnoreAll
}

// DirectMapTo just map the same name and type
// src and dst should be pointer of struct
func DirectMapTo(src, dst interface{}) {
	srcValue := reflect.ValueOf(src).Elem()
	srcType := srcValue.Type()
	dstValue := reflect.ValueOf(dst).Elem()
	// dstType := dstValue.Type()
	// log.Printf("map %v to %v", srcType.Name(), dstType.Name())
	for i := 0; i < srcType.NumField(); i++ {
		srcMemberValueField := srcValue.Field(i)
		// 过滤私有成员
		if !srcMemberValueField.CanInterface() {
			continue
		}
		// src的成员名字
		srcMemberName := srcType.Field(i).Name
		// 如果src成员在dst不存在同名，跳过
		dstMemberValueField := dstValue.FieldByName(srcMemberName)
		if !dstMemberValueField.IsValid() {
			continue
		}
		// src value member
		srcValueMember := srcMemberValueField.Interface()
		switch srcMemberValueField.Kind() {
		case reflect.Bool:
			dstMemberValueField.SetBool(srcValueMember.(bool))
		case reflect.Int:
			dstMemberValueField.SetInt(int64(srcValueMember.(int)))
		case reflect.Int8:
			dstMemberValueField.SetInt(int64(srcValueMember.(int8)))
		case reflect.Int16:
			dstMemberValueField.SetInt(int64(srcValueMember.(int16)))
		case reflect.Int32:
			dstMemberValueField.SetInt(int64(srcValueMember.(int32)))
		case reflect.Int64:
			dstMemberValueField.SetInt(srcValueMember.(int64))
		case reflect.Uint:
			dstMemberValueField.SetUint(uint64(srcValueMember.(uint)))
		case reflect.Uint8:
			dstMemberValueField.SetUint(uint64(srcValueMember.(uint8)))
		case reflect.Uint16:
			dstMemberValueField.SetUint(uint64(srcValueMember.(uint16)))
		case reflect.Uint32:
			dstMemberValueField.SetUint(uint64(srcValueMember.(uint32)))
		case reflect.Uint64:
			dstMemberValueField.SetUint(srcValueMember.(uint64))
		case reflect.Float32:
			dstMemberValueField.SetFloat(float64(srcValueMember.(float32)))
		case reflect.Float64:
			dstMemberValueField.SetFloat(srcValueMember.(float64))
		case reflect.String:
			dstMemberValueField.SetString(srcValueMember.(string))
		case reflect.Ptr:
			log.Print("ptr??")
		case reflect.Array, reflect.Slice:
			length := srcMemberValueField.Len()
			// log.Printf("found array or slice, length is %v", length)
			// dstMemberType is a slice of pointer
			dstMemberType := dstMemberValueField.Type()
			// element of slice
			dstMemberElementType := dstMemberType.Elem()
			// and maybe element of pointer
			if dstMemberElementType.Kind() == reflect.Ptr {
				dstMemberElementType = dstMemberElementType.Elem()
			}
			tmpDstMemberValueField := reflect.MakeSlice(dstMemberType, 0, length)
			for l := 0; l < length; l++ {
				srcMemberElementField := srcMemberValueField.Index(l)
				if srcMemberElementField.Kind() == reflect.Ptr {
					tmpDstMemberElementField := reflect.New(dstMemberElementType)
					DirectMapTo(srcMemberElementField.Interface(), tmpDstMemberElementField.Interface())
					tmpDstMemberValueField = reflect.Append(tmpDstMemberValueField, tmpDstMemberElementField)
				} else {
					tmpDstMemberValueField = reflect.Append(tmpDstMemberValueField, srcMemberElementField)
				}
			}
			dstMemberValueField.Set(tmpDstMemberValueField)
		case reflect.Struct:
			dstMemberValueField.Set(srcMemberValueField)
		default:
			log.Printf("don't know what to do with %v with value: %v", srcMemberName, srcMemberValueField)
		}
	}
}

func JsonMapTo(src, dst interface{}) {
	srcJson, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(srcJson), dst)
	if err != nil {
		panic(err)
	}
}
