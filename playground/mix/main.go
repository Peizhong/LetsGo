package main

// #cgo CFLAGS: -I../cpp/include
// #cgo LDFLAGS: -L../cpp -lbridge
// #include <stdlib.h>
// #include "bridge.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// export LD_LIBRARY_PATH=../cpp
func main() {
	cs := C.CString("Hello from stdio")
	v := C.Hello(11)
	v = C.SocketClient()
	fmt.Println(v)
	C.free(unsafe.Pointer(cs))
}