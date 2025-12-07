package service

import (
	"net/http"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"
	"github.com/wup364/pakku/pkg/serviceutil"
)

// HTTPService HTTP服务路由
type HTTPService struct {
	app    ipakku.Application
	http   *serviceutil.HTTPService
	config ipakku.AppConfig `@autowired:""`
}

// AsModule 作为一个模块加载
func (service *HTTPService) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "HTTPService",
		Version:     1.0,
		Description: "HTTP服务路由",
		OnReady: func(app ipakku.Application) {
			service.app = app
			service.http = serviceutil.NewHTTPService()
		},
	}
}

// SetDebug SetDebug
func (service *HTTPService) SetDebug(debug bool) {
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
	logs.Debugf("AsRouter: %T", router)

	// 自动注入依赖
	if err := service.app.Utils().AutoWired(router); nil != err {
		return err
	}
	// 自动完成配置
	if err := service.config.ScanAndAutoConfig(router); nil != err {
		return err
	}
	return service.http.AsRouter(url, router)
}

// AsController 批量注册路由, 使用RequestMapping字段作为前缀url
func (service *HTTPService) AsController(router ipakku.Controller) (err error) {
	logs.Debugf("AsController: %T", router)

	// 自动注入依赖
	if err = service.app.Utils().AutoWired(router); nil != err {
		return
	}
	// 自动完成配置
	if err = service.config.ScanAndAutoConfig(router); nil != err {
		return
	}
	//
	ctl := router.AsController()
	if err = service.http.BulkRouters(ctl.RequestMapping, ctl.ToLowerCase, ctl.HandlerFunc); nil != err {
		return
	}
	if len(ctl.FilterConfig) > 0 {
		for _, fc := range ctl.FilterConfig {
			var subPath string
			if subPath = fc.Path; len(subPath) == 0 {
				subPath = ":**"
			}
			if err = service.Filter(ctl.RequestMapping+"/"+subPath, fc.Func); nil != err {
				return
			}
		}
	}
	return
}

// Filter Filter
func (service *HTTPService) Filter(url string, fun ipakku.FilterFunc) error {
	return service.http.AddURLFilter(url, func(w http.ResponseWriter, r *http.Request) bool {
		return fun(w, r)
	})
}

// SetStaticDIR SetStaticDIR
func (service *HTTPService) SetStaticDIR(path, dir string, fun ipakku.FilterFunc) (err error) {
	if fun == nil {
		return service.http.SetStaticDIR(path, dir, nil)
	}
	return service.http.SetStaticDIR(path, dir, func(rw http.ResponseWriter, r *http.Request) bool {
		return fun(rw, r)
	})
}

// SetStaticFile SetStaticFile
func (service *HTTPService) SetStaticFile(path, file string, fun ipakku.FilterFunc) error {
	if fun == nil {
		return service.http.SetStaticFile(path, file, nil)
	}
	return service.http.SetStaticFile(path, file, func(rw http.ResponseWriter, r *http.Request) bool {
		return fun(rw, r)
	})
}
