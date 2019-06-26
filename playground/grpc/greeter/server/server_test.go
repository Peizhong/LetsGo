package main

import (
	"github.com/golang/mock/gomock"
	"github.com/peizhong/letsgo/pkg/foo"
	"github.com/peizhong/letsgo/pkg/mock_foo"
	"log"
	"reflect"
	"testing"
	"unsafe"
)

func TestGreeter_Main(t *testing.T) {
	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObj := mock_foo.NewMockFoo(ctrl)
	gomock.InOrder(
		mockObj.EXPECT().DoSomething(123, "Hello GoMock").Return(nil).Times(1),
		mockObj.EXPECT().DoSomething(22, "Hello GoMock").Return(nil).Times(1),
		mockObj.EXPECT().DoSomething(33, "Hello GoMock").Return(nil).Times(1),
	)

	func(f foo.Foo) {
		f.DoSomething(123, "Hello GoMock")
		f.DoSomething(22, "Hello GoMock")
		f.DoSomething(33, "Hello GoMock")
	}(mockObj)
}

func TestSclice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	log.Println(s)
	log.Printf(reflect.TypeOf(s).Name())
	ptr := unsafe.Pointer(&s)
	pt1 := unsafe.Pointer(&(s[1]))
	pt2 := unsafe.Pointer(&(s[2]))
	x := uintptr(pt2) - uintptr(pt1)
	log.Println(ptr)
	log.Println(x)
}
