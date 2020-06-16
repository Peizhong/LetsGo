package main

import "runtime"

func whereGoroutine() {
	runtime.Gosched()

	// g: /usr/local/go/src/runtime/runtime2.go #395
	// m: /usr/local/go/src/runtime/runtime2.go #473
	// p: /usr/local/go/src/runtime/runtime2.go #552
	// newstack: /usr/local/go/src/runtime/stack.go #918
	// sysmon: /usr/local/go/src/runtime/proc.go #4461
	// gopark: /usr/local/go/src/runtime/proc.go #287
}
