package proxy

import "fmt"

type roomService struct {
	repo Repository
}

var (
	// roomLocationKey: redis:string key(roomId)|value(endpoint)
	roomLocationKey = func(roomId string) string {
		return fmt.Sprintf("room:%s", roomId)
	}

	// endpointInfoKey: redis:score set key(endpoint)|member(roomId)|score(count)
	endpointInfoKey = func(endpoint string) string {
		return fmt.Sprintf("endpoint:%s", endpoint)
	}
)

func NewRoomSerivce() *roomService {
	return &roomService{
		repo: NewRedisRepository(),
	}
}

// Where: 查询roomId所在endpoint
func (r *roomService) Where(roomId string) (string, bool) {
	return r.repo.GetString(roomLocationKey(roomId))
}

// 加入房间，记录房间中的人数
func (r *roomService) CheckIn(endpoint, roomId string) {
	r.repo.SetString(roomLocationKey(roomId), endpoint)
	r.repo.IncrSortedSetMemberScore(endpointInfoKey(endpoint), roomId, 1)
}

// Leave: 有人离开房间。房间仍保留
func (r *roomService) Leave(endpoint, roomId string) {
	r.repo.IncrSortedSetMemberScore(endpointInfoKey(endpoint), roomId, -1)
}

// CheckOut: 删除房间信息
func (r *roomService) CheckOut(endpoint, roomId string) {
	r.repo.DeleteString(roomLocationKey(roomId))
	r.repo.SetSortedSetMemberScore(endpointInfoKey(endpoint), roomId, 0)
}

// CleanUp: 查找没有人的房间，删除
func (r *roomService) CleanUp() {
	// todo: get keys
}
