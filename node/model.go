package node

// 定义每一个模块的基础服务接口

type BasicService interface {
	Start() error
	Stop() error
}
