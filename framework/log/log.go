package log

import "fmt"

// Info log level info
func Info(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
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
