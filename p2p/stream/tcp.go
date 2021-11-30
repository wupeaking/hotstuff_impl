package stream

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/wupeakig/hotstuff_impl/conf"
	cryptogo "github.com/wupeakig/hotstuff_impl/crypto"
	"github.com/wupeakig/hotstuff_impl/node"
	"github.com/wupeakig/hotstuff_impl/p2p"
)

type TCPStream struct {
	curPeer     *p2p.Peer
	peers       *p2p.PeerBooks
	listen      net.Listener
	peerCons    map[string]io.ReadWriteCloser
	msgCallBack map[string]p2p.OnReceive
	sync.RWMutex
}

// todo::
func New(cfg *conf.Configuration) (*TCPStream, error) {
	l, err := net.Listen("tcp", cfg.ListenAddr)
	fmt.Println("listen addr: ", cfg.ListenAddr)
	if err != nil {
		return nil, err
	}
	stream := &TCPStream{
		listen:      l,
		peers:       p2p.NewPeerBooks(),
		peerCons:    map[string]io.ReadWriteCloser{},
		msgCallBack: map[string]p2p.OnReceive{},
	}

	curPeerID := cryptogo.Bytes2Hex(cfg.CurVeridier.PublickKey)
	for _, peer := range cfg.Peers {
		if curPeerID == peer.ID {
			stream.curPeer = &p2p.Peer{
				ID:      peer.ID,
				Address: peer.Address,
			}
		}
		stream.peers.AddPeer(&p2p.Peer{
			ID:      peer.ID,
			Address: peer.Address,
		})
	}
	return stream, nil
}

// impl node.BasicService

var _ node.BasicService = &TCPStream{}

func (ts *TCPStream) Start() error {
	go func() {
		for {
			conn, err := ts.listen.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			go ts.handleConnect(conn)

		}
	}()

	return nil
}

func (ts *TCPStream) Stop() error {
	// TODO::
	ts.listen.Close()
	return nil
}

func (ts *TCPStream) handleConnect(con net.Conn) {
	// 如果已经存在连接 则关闭此连接

	msg, err := ts.unpackageData(bufio.NewReader(con))
	if err != nil {
		fmt.Printf("读取消息出错, err: %s\n", err.Error())
		con.Close()
		return
	}

	ts.RLock()
	msgCB := ts.msgCallBack[msg.ModelID]
	ts.RUnlock()
	if msgCB == nil {
		fmt.Printf("未知的model id : %s, peer: %s, address: %s", msg.ModelID, ts.curPeer.ID, ts.curPeer.Address)
		con.Close()
		return
	}
	go msgCB(msg.ModelID, msg.Msg, &p2p.Peer{ID: msg.PeerID})

	ts.Lock()
	_, exist := ts.peerCons[msg.PeerID]
	if exist {
		con.Close()
		ts.Unlock()
		return
	}
	ts.peerCons[msg.PeerID] = con
	ts.Unlock()

	// 一直读取消息
	for {
		msg, err := ts.unpackageData(bufio.NewReader(con))
		if err != nil {
			fmt.Printf("读取消息出错, err: %s\n", err.Error())
			con.Close()
			return
		}
		ts.RLock()
		msgCB := ts.msgCallBack[msg.ModelID]
		ts.RUnlock()
		if msgCB == nil {
			fmt.Println("未知的model id")
			continue
		}
		go msgCB(msg.ModelID, msg.Msg, &p2p.Peer{ID: msg.PeerID})
	}
}

var _ p2p.SwitcherI = &TCPStream{}

// 向所有的节点广播消息
func (ts *TCPStream) Broadcast(modelID string, msg *p2p.BroadcastMsg) error {
	peers := ts.peers.AllPeers()
	for _, peer := range peers {
		err := ts.BroadcastToPeer(modelID, msg, peer)
		if err != nil {
			fmt.Printf("Broadcast err: %v\n", err)
		}
	}
	return nil
}

// 广播到指定的peer
func (ts *TCPStream) BroadcastToPeer(modelID string, msg *p2p.BroadcastMsg, p *p2p.Peer) error {
	peer := ts.peers.FindPeer(p.ID)
	if peer == nil {
		return nil
	}
	if peer.ID == ts.curPeer.ID {
		println("本机节点 不广播......", ts.curPeer.ID)
		return nil
	}
	binMsg, err := ts.packageData(msg)
	if err != nil {
		return err
	}

	// 首先尝试获取已有连接
	ts.RLock()
	con := ts.peerCons[peer.ID]
	ts.RUnlock()
	if con != nil {
		_, err := io.Copy(con, bytes.NewBuffer(binMsg))
		if err != nil {
			fmt.Printf("广播出错 err: %s\n", err.Error())
			ts.Lock()
			delete(ts.peerCons, msg.PeerID)
			ts.Unlock()
			con.Close()
		}
		return nil
	}
	//说明没有建立连接
	tcpAddr, err := net.ResolveTCPAddr("tcp", peer.Address)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	_, err = io.Copy(conn, bytes.NewBuffer(binMsg))
	if err != nil {
		return err
	}
	ts.Lock()
	ts.peerCons[peer.ID] = conn
	ts.Unlock()
	return nil
}

// 广播 除了指定的peer
func (ts *TCPStream) BroadcastExceptPeer(modelID string, msg *p2p.BroadcastMsg, p *p2p.Peer) error {
	peers := ts.peers.AllPeers()
	for _, peer := range peers {
		if peer.ID == p.ID {
			continue
		}
		ts.BroadcastToPeer(modelID, msg, peer)
	}
	return nil
}

// 移除某个peer
func (ts *TCPStream) RemovePeer(p *p2p.Peer) error {
	ts.peers.RemovePeer(p.ID)
	return nil
}

func (ts *TCPStream) RegisterOnReceive(modelID string, callBack p2p.OnReceive) error {
	ts.Lock()
	ts.msgCallBack[modelID] = callBack
	ts.Unlock()
	return nil
}

// 返回所有存在的peers
func (ts *TCPStream) Peers() ([]*p2p.Peer, error) {
	return ts.peers.AllPeers(), nil
}
