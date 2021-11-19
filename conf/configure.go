package conf

import (
	"os"

	cryptogo "github.com/wupeakig/hotstuff_impl/crypto"
	"github.com/wupeakig/hotstuff_impl/model"
	"github.com/wupeakig/hotstuff_impl/p2p"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Veridiers   []*model.Verifier
	Peers       []p2p.Peer
	CurVeridier *model.Verifier
	CurSeqNum   int
	ListenAddr  string
}

type Cfg struct {
	Veridiers []struct {
		PublickKey string `yaml:"public_key"`
		PrivateKey string `yaml:"private_key"`
		Address    string `yaml:"address"`
	} `yaml:"veridiers"`
	PrivateKey string `yaml:"private_key"`
	ListenAddr string `yaml:"listen_addr"`
}

func NewConfiguration() *Configuration {
	fd, err := os.Open("./conf/config.yaml")
	if err != nil {
		panic(err)
	}
	cfg := Cfg{}
	err = yaml.NewDecoder(fd).Decode(&cfg)
	if err != nil {
		panic(err)
	}
	pri, err := cryptogo.LoadPrivateKey(cfg.PrivateKey)
	if err != nil {
		panic(err)
	}
	pubKey := cryptogo.PublicKey2Hex(&pri.PublicKey)

	retCfg := &Configuration{
		Veridiers: make([]*model.Verifier, 0, len(cfg.Veridiers)),
		Peers:     make([]p2p.Peer, 0, len(cfg.Veridiers)),
		CurSeqNum: -1,
	}
	for i, v := range cfg.Veridiers {
		pubBytes, err := cryptogo.Hex2Bytes(v.PublickKey)
		if err != nil {
			panic(err)
		}
		priBytes, err := cryptogo.Hex2Bytes(v.PrivateKey)
		if err != nil {
			panic(err)
		}
		if v.Address == "" {
			panic("address is empty")
		}
		retCfg.Veridiers = append(retCfg.Veridiers, &model.Verifier{
			PublickKey: pubBytes,
			PrivateKey: priBytes,
			SeqNum:     int32(i),
		})
		retCfg.Peers = append(retCfg.Peers, p2p.Peer{
			ID:      v.PublickKey,
			Address: v.Address,
		})

		if v.PublickKey == pubKey {
			retCfg.CurVeridier = &model.Verifier{
				PublickKey: pubBytes,
				PrivateKey: priBytes,
				SeqNum:     int32(i),
			}
		}
	}

	if retCfg.CurVeridier == nil {
		panic("未指定当前验证者")
	}
	retCfg.ListenAddr = cfg.ListenAddr

	return retCfg
}
