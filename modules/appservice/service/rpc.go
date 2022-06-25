package service

import (
	"net/rpc"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
)

// RPCService RPC服务路由
type RPCService struct {
	isdebug bool
	mctx    ipakku.Loader
	rpcs    *rpc.Server
}

// AsModule 作为一个模块加载
func (rpcs *RPCService) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "RPCService",
		Version:     1.0,
		Description: "RPC服务路由",
		OnReady: func(mctx ipakku.Loader) {
			rpcs.mctx = mctx
			rpcs.rpcs = rpc.NewServer()
		},
		OnSetup:  func() {},
		OnUpdate: func(cv float64) {},
		OnInit: func() {
		},
	}
}

// SetDebug SetDebug
func (rpcs *RPCService) SetDebug(debug bool) {
	rpcs.isdebug = debug
}

// GetRPCService GetRPCService
func (rpcs *RPCService) GetRPCService() *rpc.Server {
	return rpcs.rpcs
}

// RegisteRPC RegisteRPC
func (rpcs *RPCService) RegisteRPC(rcvr interface{}) error {
	// 自动注入依赖
	if err := rpcs.mctx.AutoWired(rcvr); nil != err {
		return err
	}
	if rpcs.isdebug {
		logs.Debugf("AddRPCService: %T\r\n", rcvr)
	}
	return rpcs.rpcs.Register(rcvr)
}
