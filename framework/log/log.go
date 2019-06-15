package log

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

// Info log level info
func Info(args ...interface{}) {
	fmt.Println(spew.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	fmt.Println(spew.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	fmt.Println(spew.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	fmt.Println(spew.Sprintf(format, args...))
}

func WithField(str ...string) *Entry {
	return &Entry{}
}

type Entry struct {
}

func (e *Entry) Info(format string, args ...interface{}) {

}

func (e *Entry) Error(format string, args ...interface{}) {

}
