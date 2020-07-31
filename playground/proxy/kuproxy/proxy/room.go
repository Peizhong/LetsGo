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

func NewRoomSerivce(config Config) *roomService {
	return &roomService{
		repo: NewRedisRepository(config),
	}
}

// Where: 查询roomId所在endpoint
func (r *roomService) Where(roomId string) (string, bool) {
	return r.repo.GetString(roomLocationKey(roomId))
}

// 申请房间，记录房间中的人数
func (r *roomService) CheckIn(endpoint, roomId string) (err error) {
	err = r.repo.SetString(roomLocationKey(roomId), endpoint)
	if err != nil {
		return
	}
	err = r.repo.SetSortedSetMemberScore(endpointInfoKey(endpoint), roomId, 1)
	return
}

// 加入房间
func (r *roomService) Enter(endpoint, roomId string) error {
	_, err := r.repo.IncrSortedSetMemberScore(endpointInfoKey(endpoint), roomId, 1)
	return err
}

// Leave: 有人离开房间。如果走完了，删除
func (r *roomService) Leave(endpoint, roomId string) {
	// 降低一个
	remain, _ := r.repo.IncrSortedSetMemberScore(endpointInfoKey(endpoint), roomId, -1)
	if remain <= 0 && remain != RepoErrorVal {
		// todo: socket.io建立握手时，向发送服务器发送三次请求。一个连接断开后，不能马上删除对应的endpoint信息
		// r.CheckOut(endpoint, roomId)
	}
}

// CheckOut: 最后一个连接端断开，删除房间信息
func (r *roomService) CheckOut(endpoint, roomId string) {
	r.repo.DeleteString(roomLocationKey(roomId))
	r.repo.RemoveSortedSetMember(endpointInfoKey(endpoint), roomId)
}

func (r *roomService) Status(endpoint string) []string {
	var res []string
	for room, num := range r.repo.GetSortedSetMembersWithScore(endpointInfoKey(endpoint)) {
		res = append(res, fmt.Sprintf("%s:%v", room, num))
	}
	return res
}
