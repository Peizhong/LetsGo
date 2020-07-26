package main

import (
	"github.com/golang/mock/gomock"
	mock_main "github.com/peizhong/letsgo/playground/proxy/kuproxy/mock"
)

// select upstream pod

type selector interface {
	// tellme 根据标识，获得对应endpoints
	TellMe(string) string
}

const (
	upstream = "192.168.3.143:3000"
)

func getSelector() (selector, error) {
	ctrl := gomock.NewController(nil)
	// defer ctrl.Finish()
	mock := mock_main.NewMockselector(ctrl)
	mock.EXPECT().TellMe(gomock.Any()).Return(upstream)
	return mock, nil
}
