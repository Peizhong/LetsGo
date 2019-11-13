package main

// #include <stdio.h>
// #include <stdlib.h>
// #include "bridge.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	// load c++
	// wrap c++ code with a c interface
	cs := C.CString("Hello from stdio")
	v := C.Hello(11)
	fmt.Println(v)
	C.free(unsafe.Pointer(cs))
}