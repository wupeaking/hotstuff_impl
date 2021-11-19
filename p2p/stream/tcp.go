package stream

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/wupeakig/hotstuff_impl/conf"
	"github.com/wupeakig/hotstuff_impl/node"
	"github.com/wupeakig/hotstuff_impl/p2p"
)

type TCPStream struct {
	peers    *p2p.PeerBooks
	listen   net.Listener
	peerCons map[string]io.ReadWriteCloser
	sync.RWMutex
}

// todo::
func New(cfg *conf.Configuration) (*TCPStream, error) {
	l, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		return nil, err
	}
	stream := &TCPStream{
		listen: l,
		peers:  p2p.NewPeerBooks(),
	}
	for _, peer := range cfg.Peers {
		stream.peers.AddPeer(&peer)
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

	return nil
}

func (ts *TCPStream) handleConnect(con net.Conn) error {
	return nil
}

// Close the listener when the application closes.
// defer l.Close()
// fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
// for {
// 	// Listen for an incoming connection.
// 	conn, err := l.Accept()
// 	if err != nil {
// 		fmt.Println("Error accepting: ", err.Error())
// 		os.Exit(1)
// 	}
// 	// Handle connections in a new goroutine.
// 	go handleRequest(conn)
// }
