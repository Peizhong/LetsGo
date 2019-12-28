package app

import (
	"encoding/json"
	"github.com/hashicorp/memberlist"
	"log"
)

type delegate struct {
	storage *storage
}

// NodeMeta is used to retrieve meta-data about the current node
// when broadcasting an alive message. It's length is limited to
// the given byte size. This metadata is available in the Node structure.
func (d *delegate) NodeMeta(limit int) []byte {
	return []byte(d.storage.Name)
}

// NotifyMsg is called when a user-data message is received.
// Care should be taken that this method does not block, since doing
// so would block the entire UDP packet receive loop. Additionally, the byte
// slice may be modified after the call returns, so it should be copied if needed
func (d *delegate) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}
	switch b[0] {
	case 'd': // data
		d.storage.UnpackNotify(b[1:])
	}
}

// GetBroadcasts is called when user data messages can be broadcast.
// It can return a list of buffers to send. Each buffer should assume an
// overhead as provided with a limit on the total byte size allowed.
// The total byte size of the resulting data to send must not exceed
// the limit. Care should be taken that this method does not block,
// since doing so would block the entire UDP packet receive loop.
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	// QueueBroadcast的数据再取出来
	// b := d.storage.broadcasts.GetBroadcasts(overhead, limit)
	b := [][]byte{}
	return b
}

// LocalState is used for a TCP Push/Pull. This is sent to
// the remote side in addition to the membership information. Any
// data can be sent here. See MergeRemoteState as well. The `join`
// boolean indicates this is for a join instead of a push/pull.
func (d *delegate) LocalState(join bool) []byte {
	d.storage.lock.RLock()
	m := d.storage.data
	d.storage.lock.RUnlock()
	b, _ := json.Marshal(m)
	return b
}

// MergeRemoteState is invoked after a TCP Push/Pull. This is the
// state received from the remote side and is the result of the
// remote side's LocalState call. The 'join'
// boolean indicates this is for a join instead of a push/pull.
func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) < 1 {
		return
	}
	log.Println("MergeRemoteState join:", join)
	if !join {
		return
	}
	var data dataStore
	if err := json.Unmarshal(buf, &data); err != nil {
		return
	}
	d.storage.lock.Lock()
	for k, v := range data {
		d.storage.data[k] = v
	}
	d.storage.lock.Unlock()
}

type alive struct{}

// NotifyAlive can filter out alive messages based on custom logic
func (a *alive) NotifyAlive(peer *memberlist.Node) error {
	// NotifyAlive is invoked when a message about a live
	// node is received from the network.  Returning a non-nil
	// error prevents the node from being considered a peer.
	log.Println("alive", peer.Name)
	return nil
}
