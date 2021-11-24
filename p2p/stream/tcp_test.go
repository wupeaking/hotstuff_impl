package stream_test

import (
	"bytes"
	"log"
	"time"

	"github.com/wupeakig/hotstuff_impl/conf"
	"github.com/wupeakig/hotstuff_impl/p2p"
	"github.com/wupeakig/hotstuff_impl/p2p/stream"
)

func ExampleTCPStream() {
	cfg1 := `
veridiers:
  - public_key: 0x7f096e2682f086736c7b3c3da10fc140c32dca71db48aca5f19bd7d825b0c7daca54978823334ef9f35ccb32dcfe44edc2015e889de821d3e5b19a65a110f877
    address: "127.0.0.1:11080"
  - public_key: 0x8f2696ae0dfd21a2736ab100d62a1300901389277e4507ff77d98ae033973fa3c2c216518f34491c268373ce7be9aa641e45fc7b0822ad00eb06f6a4c08da70c
    address: "127.0.0.1:21080"
  - public_key: 0x7032f77b9fb66665a49b2d6e1ddd2c0a03c62aa9a2026e299a192fd1b765ec539361e65c85edacf799bffae2716053e15b78498ff75b42a057d73e719be7e4fb
    address: "127.0.0.1:31080"
listen_addr: 127.0.0.1:11080
private_key: 0x325cf070c843a6de85dc8ed92b975544a0097a8ab1436bb47d83558e296c5b9f
`

	cfg2 := `
veridiers:
  - public_key: 0x7f096e2682f086736c7b3c3da10fc140c32dca71db48aca5f19bd7d825b0c7daca54978823334ef9f35ccb32dcfe44edc2015e889de821d3e5b19a65a110f877
    address: "127.0.0.1:11080"
  - public_key: 0x8f2696ae0dfd21a2736ab100d62a1300901389277e4507ff77d98ae033973fa3c2c216518f34491c268373ce7be9aa641e45fc7b0822ad00eb06f6a4c08da70c
    address: "127.0.0.1:21080"
  - public_key: 0x7032f77b9fb66665a49b2d6e1ddd2c0a03c62aa9a2026e299a192fd1b765ec539361e65c85edacf799bffae2716053e15b78498ff75b42a057d73e719be7e4fb
    address: "127.0.0.1:31080"
listen_addr: 127.0.0.1:21080
private_key: 0x3bf83bae9d0e5917687ee5039f0bdd13a91bff6a179688d8bea7ff3b69166e26
`
	tcpNode1, err := stream.New(conf.NewConfiguration(cfg1))
	if err != nil {
		log.Fatal(err)
	}
	if err := tcpNode1.Start(); err != nil {
		log.Fatal(err)
	}

	modelID := "test"
	wait := false
	msgCB := func(model string, msgBytes []byte, p *p2p.Peer) {
		if model != modelID {
			panic("model id err")
		}
		if !bytes.Equal(msgBytes, []byte{1, 2, 3, 4}) {
			panic("msg err")
		}
		wait = true
	}

	tcpNode1.RegisterOnReceive(modelID, msgCB)

	tcpNode2, err := stream.New(conf.NewConfiguration(cfg2))
	if err != nil {
		log.Fatal(err)
	}
	if err := tcpNode2.Start(); err != nil {
		log.Fatal(err)
	}

	err = tcpNode2.Broadcast(modelID,
		&p2p.BroadcastMsg{ModelID: modelID, Msg: []byte{1, 2, 3, 4}})
	if err != nil {
		panic(err)
	}

	allpeers, _ := tcpNode2.Peers()
	for i := range allpeers {
		println("peers: ", allpeers[i].Address)
	}

	for {
		time.Sleep(1 * time.Second)
		println("wait-----")
		if wait {
			break
		}
	}
	// Output:
}
