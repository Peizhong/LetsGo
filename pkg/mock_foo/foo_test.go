package mock_foo

import (
	"github.com/golang/mock/gomock"
	"github.com/peizhong/letsgo/pkg/foo"
	"testing"
)

func TestSUT(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := NewMockFoo(ctrl)

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	m.
		EXPECT().
		Bar(gomock.Eq(99)).
		Return(101).
		AnyTimes()

	foo.SUT(m)
}
