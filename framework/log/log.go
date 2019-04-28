package log

import "fmt"

// Info log level info
func Info(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
