package appservice

import (
	"net"
	"net/http"
	"time"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/modules/appservice/service"
	"github.com/wup364/pakku/utils/logs"
)

// AppService HTTP服务
type AppService struct {
	service.RPCService
	service.HTTPService
	conf ipakku.AppConfig `@autowired:""`
}

// AsModule 作为一个模块加载
func (service *AppService) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Version:     1.0,
		Description: "AppService module",
		OnReady: func(app ipakku.Application) {
			if nil != service.HTTPService.AsModule().OnReady {
				service.HTTPService.AsModule().OnReady(app)
			}
			if nil != service.RPCService.AsModule().OnReady {
				service.RPCService.AsModule().OnReady(app)
			}
		},
	}
}

// StartHTTP 启动服务
func (service *AppService) StartHTTP(serviceCfg ipakku.HTTPServiceConfig) {
	service.HTTPService.SetDebug(serviceCfg.Debug)

	s := serviceCfg.Server
	if nil == s {
		s = &http.Server{
			WriteTimeout:   time.Duration(service.conf.GetConfig(ipakku.CONFKEY_WRITETIMEOUTSECOND).ToInt64(0)) * time.Second,
			ReadTimeout:    time.Duration(service.conf.GetConfig(ipakku.CONFKEY_READTIMEOUTSECOND).ToInt64(0)) * time.Second,
			MaxHeaderBytes: service.conf.GetConfig(ipakku.CONFKEY_MAXHEADERBYTES).ToInt(1 << 20),
		}
	}

	if nil == s.ErrorLog {
		s.ErrorLog = logs.ErrorLogger()
	}

	s.Handler = service.GetRouter()
	if len(serviceCfg.ListenAddr) > 0 {
		s.Addr = serviceCfg.ListenAddr
	}

	var err error
	if len(serviceCfg.CertFile) == 0 || len(serviceCfg.KeyFile) == 0 {
		logs.Infoln("Server(HTTP) listened in: " + s.Addr)
		err = s.ListenAndServe()
	} else {
		logs.Infoln("Server(HTTPS) listened in: " + s.Addr)
		err = s.ListenAndServeTLS(serviceCfg.CertFile, serviceCfg.KeyFile)
	}

	if nil != err {
		logs.Panicln(err)
	}
}

// StartRPC 启动服务
func (service *AppService) StartRPC(serviceCfg ipakku.RPCServiceConfig) {
	service.RPCService.SetDebug(serviceCfg.Debug)

	if len(serviceCfg.Network) == 0 {
		serviceCfg.Network = "tcp"
	}

	l := serviceCfg.Listener
	if nil == l {
		var err error
		if l, err = net.Listen(serviceCfg.Network, serviceCfg.ListenAddr); nil != err {
			logs.Panicln(err)
		}
	}
	s := service.RPCService.GetRPCService()
	logs.Infoln("Server(RPC) listened in: " + serviceCfg.ListenAddr)
	if err := http.Serve(l, s); err != nil {
		logs.Panicln(err)
	}
}
