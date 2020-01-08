package main

// #cgo CFLAGS: -I../cpp/include
// #cgo LDFLAGS: -L../cpp -lbridge
// #include <stdio.h>
// #include <stdlib.h>
// #include "bridge.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type goClient struct{
	client C.WrapClient
}

func NewGoClient() goClient{
	var ret goClient
	ret.client = C.WrapClientInit()
	return ret
}

func (g goClient) SendMessage(message string) string {
	msg := C.CString(message)
	defer C.free(unsafe.Pointer(msg))
	r := C.WrapClientSendMessage(unsafe.Pointer(g.client),msg)
	resp := C.GoString(r)
	fmt.Println(resp)
	return resp
}

// export LD_LIBRARY_PATH=../cpp
func main() {
	client := NewGoClient()
	client.SendMessage("hello")
}