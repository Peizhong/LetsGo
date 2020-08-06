package service

import (
	"github.com/peizhong/letsgo/playground/ast/model"
	"github.com/peizhong/letsgo/playground/ast/runtime"
)

type DemoService struct {
	Runtime runtime.DemoRuntime
}

func (d *DemoService) Get(r *model.DemoRequest) (*model.DemoResponse, error) {
	return &model.DemoResponse{}, nil
}

const template = `
package test_service

import (
	// copy from source
)

func NewXXSourceStruct() *XXSourceStruct {
	// 可以使用Mock
}

func TestXXXSourceStructXXFunc(t *testing.T){
	service := NewXXSourceStruct()
	xxresp,err:=service.XXFunc(xxrequest)
	assert.NoError(t,err)
}
`
