package service

import (
	"net/rpc"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
)

// RPCService RPC服务路由
type RPCService struct {
	isdebug bool
	rpcs    *rpc.Server
	app     ipakku.Application
	config  ipakku.AppConfig `@autowired:""`
}

// AsModule 作为一个模块加载
func (rpcs *RPCService) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "RPCService",
		Version:     1.0,
		Description: "RPC服务路由",
		OnReady: func(app ipakku.Application) {
			rpcs.app = app
			rpcs.rpcs = rpc.NewServer()
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
	logs.Debugf("AddRPCService: %T\r\n", rcvr)

	// 自动注入依赖
	if err := rpcs.app.Utils().AutoWired(rcvr); nil != err {
		return err
	}
	// 自动完成配置
	if err := rpcs.config.ScanAndAutoConfig(rcvr); nil != err {
		return err
	}
	return rpcs.rpcs.Register(rcvr)
}
