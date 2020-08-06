package service

import (
	"context"

	"github.com/peizhong/letsgo/playground/ast/model"
	"github.com/peizhong/letsgo/playground/ast/runtime"
)

type DemoService struct {
	Runtime runtime.DemoRuntime
}

func (d *DemoService) Get(ctx context.Context, r *model.DemoRequest) (*model.DemoResponse, error) {
	return &model.DemoResponse{}, nil
}

func (d *DemoService) Gets(ctx context.Context, r *model.DemoRequest) (*model.DemoResponse, int64, error) {
	return &model.DemoResponse{}, 0, nil
}

func (d *DemoService) Update(ctx context.Context, r *model.DemoRequest) (bool, error) {
	return true, nil
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
