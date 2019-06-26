package utils

import (
	"errors"
	"html/template"
	"log"
	"os"
	"reflect"
)

// FuncParamMatch check if args fit func
func FuncParamMatch(f interface{}, args ...interface{}) error {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return errors.New("not func")
	}
	if t.NumIn() != len(args) {
		return errors.New("in params count not match")
	}
	todo := make([]reflect.Value, len(args))
	for i, a := range args {
		p := t.In(i)
		if reflect.TypeOf(a) != p {
			return errors.New("not match for " + p.String())
		}
		todo[i] = reflect.ValueOf(a)
	}
	v := reflect.ValueOf(f)
	res := v.Call(todo)
	for _, r := range res {
		log.Println(r.Interface())
	}
	return nil
}

func Template2() {
	const test = `Hello, [[.Name]]`
	tpl, err := template.New("template-name").Delims("[[", "]]").Parse(test)
	if err == nil {
		err = tpl.Execute(os.Stdout, map[string]interface{}{
			"Name": "wpz",
		})
	}
}
