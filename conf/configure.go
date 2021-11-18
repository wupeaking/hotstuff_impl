package conf

import (
	"os"

	cryptogo "github.com/wupeakig/hotstuff_impl/crypto"
	"github.com/wupeakig/hotstuff_impl/model"
	"github.com/wupeakig/hotstuff_impl/p2p"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Veridiers []*model.Verifier `yaml:"veridiers"`
	Peers     []p2p.Peer        ``
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
	return &Configuration{}
}
