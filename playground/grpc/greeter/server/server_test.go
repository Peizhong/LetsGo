package main

import (
	"github.com/golang/mock/gomock"
	"github.com/peizhong/letsgo/pkg/foo"
	"github.com/peizhong/letsgo/pkg/mock_foo"
	"testing"
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
