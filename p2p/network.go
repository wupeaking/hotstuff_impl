package p2p

import (
	"sync"
	// "github.com/wupeakig/hotstuff_impl/model"
)

type SwitcherI interface {
	// 向所有的节点广播消息
	Broadcast(modelID string, msg *BroadcastMsg) error
	// 广播到指定的peer
	BroadcastToPeer(modelID string, msg *BroadcastMsg, p *Peer) error
	// 广播 除了指定的peer
	BroadcastExceptPeer(modelID string, msg *BroadcastMsg, p *Peer) error
	// 移除某个peer
	RemovePeer(p *Peer) error
	RegisterOnReceive(modelID string, callBack OnReceive) error
	// 返回所有存在的peers
	Peers() ([]*Peer, error)
	// Recv() <-chan interface{}
}

type BroadcastMsg struct {
	ModelID string `json:"model_id"`
	Msg     []byte `json:"msg"`
}

// OnReceive 注册接收消息回调
type OnReceive func(modelID string, msgBytes []byte, p *Peer)

type Peer struct {
	ID      string // 定义peerid  每个peerid应该是唯一的
	Address string // 地址
}

type PeerBooks struct {
	sync.RWMutex
	sets map[string]*Peer
}

func NewPeerBooks() *PeerBooks {
	return &PeerBooks{
		sets: make(map[string]*Peer),
	}
}

func (pb *PeerBooks) AddPeer(p *Peer) {
	if p == nil {
		return
	}
	pb.Lock()
	pb.sets[p.ID] = p
	pb.Unlock()
}

func (pb *PeerBooks) FindPeer(id string) *Peer {
	pb.RLock()
	defer pb.RUnlock()
	v := pb.sets[id]
	return v
}

func (pb *PeerBooks) RemovePeer(id string) {
	pb.Lock()
	defer pb.Unlock()
	delete(pb.sets, id)
}
