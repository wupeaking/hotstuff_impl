package conf

import (
	"github.com/wupeakig/hotstuff_impl/model"
	"github.com/wupeakig/hotstuff_impl/p2p"
)

type Configuration struct {
	Veridiers []*model.Verifier
	Peers     []p2p.Peer
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}
