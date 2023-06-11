package service

import (
	"net/http"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/serviceutil"
)

// HTTPService HTTP服务路由
type HTTPService struct {
	isdebug bool
	mctx    ipakku.Loader
	config  ipakku.AppConfig
	http    *serviceutil.HTTPService
}

// AsModule 作为一个模块加载
func (service *HTTPService) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "HTTPService",
		Version:     1.0,
		Description: "HTTP服务路由",
		OnReady: func(mctx ipakku.Loader) {
			service.mctx = mctx
			service.http = serviceutil.NewHTTPService()
			if err := mctx.GetModules(&service.config); nil != err {
				logs.Panicln(err)
			}
		},
	}
}

// SetDebug SetDebug
func (service *HTTPService) SetDebug(debug bool) {
	service.isdebug = debug
	service.http.SetDebug(debug)
}

// GetRouter GetRouter
func (service *HTTPService) GetRouter() *serviceutil.ServiceRouter {
	return service.http.GetRouter()
}

// Get Get
func (service *HTTPService) Get(url string, fun ipakku.HandlerFunc) error {
	return service.http.Get(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// Post Post
func (service *HTTPService) Post(url string, fun ipakku.HandlerFunc) error {
	return service.http.Post(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// Put Put
func (service *HTTPService) Put(url string, fun ipakku.HandlerFunc) error {
	return service.http.Put(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// Patch Patch
func (service *HTTPService) Patch(url string, fun ipakku.HandlerFunc) error {
	return service.http.Patch(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// Head Head
func (service *HTTPService) Head(url string, fun ipakku.HandlerFunc) error {
	return service.http.Head(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// Options Options
func (service *HTTPService) Options(url string, fun ipakku.HandlerFunc) error {
	return service.http.Options(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// Delete Delete
func (service *HTTPService) Delete(url string, fun ipakku.HandlerFunc) error {
	return service.http.Delete(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// Any Any
func (service *HTTPService) Any(url string, fun ipakku.HandlerFunc) error {
	return service.http.Any(url, func(w http.ResponseWriter, r *http.Request) {
		fun(w, r)
	})
}

// AsRouter 批量注册路由, 可以再指定一个前缀url
func (service *HTTPService) AsRouter(url string, router ipakku.Router) error {
	if service.isdebug {
		logs.Debugf("AsRouter: %T\r\n", router)
	}
	// 自动注入依赖
	if err := service.mctx.AutoWired(router); nil != err {
		return err
	}
	// 自动完成配置
	if err := service.config.ScanAndAutoConfig(router); nil != err {
		return err
	}
	return service.http.AsRouter(url, router)
}

// AsController 批量注册路由, 使用RequestMapping字段作为前缀url
func (service *HTTPService) AsController(router ipakku.Controller) error {
	if service.isdebug {
		logs.Debugf("AsController: %T\r\n", router)
	}
	// 自动注入依赖
	if err := service.mctx.AutoWired(router); nil != err {
		return err
	}
	// 自动完成配置
	if err := service.config.ScanAndAutoConfig(router); nil != err {
		return err
	}
	ctl := router.AsController()
	if err := service.http.BulkRouters(ctl.RequestMapping, ctl.ToLowerCase, ctl.HandlerFunc); nil == err {
		if len(ctl.FilterConfig.FilterFunc) > 0 {
			for _, fc := range ctl.FilterConfig.FilterFunc {
				if val, ok := fc[1].(ipakku.FilterFunc); ok {
					if err = service.Filter(ctl.RequestMapping+"/"+fc[0].(string), val); nil != err {
						return err
					}
				} else {
					if err = service.Filter(ctl.RequestMapping+"/"+fc[0].(string), func(rw http.ResponseWriter, r *http.Request) bool {
						return fc[1].(func(http.ResponseWriter, *http.Request) bool)(rw, r)
					}); nil != err {
						return err
					}
				}
			}
		}
	} else {
		return err
	}
	return nil
}

// Filter Filter
func (service *HTTPService) Filter(url string, fun ipakku.FilterFunc) error {
	return service.http.AddURLFilter(url, func(w http.ResponseWriter, r *http.Request) bool {
		return fun(w, r)
	})
}

// SetStaticDIR SetStaticDIR
func (service *HTTPService) SetStaticDIR(path, dir string, fun ipakku.FilterFunc) (err error) {
	return service.http.SetStaticDIR(path, dir, func(rw http.ResponseWriter, r *http.Request) bool {
		if nil == fun {
			return true
		}
		return fun(rw, r)
	})
}

// SetStaticFile SetStaticFile
func (service *HTTPService) SetStaticFile(path, file string, fun ipakku.FilterFunc) error {
	return service.http.SetStaticFile(path, file, func(rw http.ResponseWriter, r *http.Request) bool {
		if nil == fun {
			return true
		}
		return fun(rw, r)
	})
}
