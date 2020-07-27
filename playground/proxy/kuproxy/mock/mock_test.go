package mock_main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	ctrl := gomock.NewController(t)

	upStream := []string{"192.168.3.143:3000"}
	mockDiscovery := NewMockdiscovery(ctrl)
	mockDiscovery.EXPECT().Endpoints().Return(upStream)
	endpoints := mockDiscovery.Endpoints()
	assert.NotEmpty(t, endpoints)

	roomId := "abc"
	mockSelector := NewMockselector(ctrl)
	mockSelector.EXPECT().TellMe(gomock.AssignableToTypeOf(roomId)).Return(endpoints[0])
	endpoint := mockSelector.TellMe(roomId)
	assert.NotZero(t, endpoint)

	mockRoom := NewMockroom(ctrl)
	mockRoom.EXPECT().Join(gomock.AssignableToTypeOf(roomId), gomock.AssignableToTypeOf(endpoint))
	mockRoom.EXPECT().Leave(gomock.AssignableToTypeOf(roomId), gomock.AssignableToTypeOf(endpoint))
	mockRoom.Join(roomId, endpoint)
	mockRoom.Leave(roomId, endpoint)
}
