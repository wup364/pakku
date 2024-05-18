// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package serviceutil

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/reflectutil"
)

// NewHTTPService 新建一个HTTP服务
func NewHTTPService() *HTTPService {
	return &HTTPService{}
}

// StartHTTPConf 启动配置
type StartHTTPConf struct {
	CertFile   string
	KeyFile    string
	ListenAddr string
	Server     *http.Server
}

// AsRouter 批量注册服务路径
type AsRouter interface {
	AsRouter() RouterConfig
}

// RouterConfig 批量注册器
type RouterConfig struct {
	// ToLowerCase  是否需要url转小写, 在未指定url(使用函数名字作为url的一部分)的情况下生效
	ToLowerCase bool

	// HandlerFunc 需要注册的函数 [{"Method(GET|POST...)", "HandlerFunc function"}, {"Method(GET|POST...)", "指定的url(可选参数)", "HandlerFunc function"}]
	HandlerFunc [][]interface{}
}

// HTTPService HTTP服务
type HTTPService struct {
	ServiceRouter
}

// GetRouter 获取路由实例, 可作为http的Server.Handler实例
func (service *HTTPService) GetRouter() *ServiceRouter {
	return &service.ServiceRouter
}

// Get Get
func (service *HTTPService) Get(url string, fun HandlerFunc) error {
	return service.AddHandler(http.MethodGet, url, fun)
}

// Post Post
func (service *HTTPService) Post(url string, fun HandlerFunc) error {
	return service.AddHandler(http.MethodPost, url, fun)
}

// Put Put
func (service *HTTPService) Put(url string, fun HandlerFunc) error {
	return service.AddHandler(http.MethodPut, url, fun)
}

// Patch Patch
func (service *HTTPService) Patch(url string, fun HandlerFunc) error {
	return service.AddHandler(http.MethodPatch, url, fun)
}

// Head Head
func (service *HTTPService) Head(url string, fun HandlerFunc) error {
	return service.AddHandler(http.MethodHead, url, fun)
}

// Options Options
func (service *HTTPService) Options(url string, fun HandlerFunc) error {
	return service.AddHandler(http.MethodOptions, url, fun)
}

// Delete Delete
func (service *HTTPService) Delete(url string, fun HandlerFunc) error {
	return service.AddHandler(http.MethodDelete, url, fun)
}

// Any Any
func (service *HTTPService) Any(url string, fun HandlerFunc) error {
	return service.AddHandler("ANY", url, fun)
}

// BulkRouters 批量注册路由, 指定一个前缀url
// routers: 需要注册的处理函数 [{"Method(GET|POST...)", "HandlerFunc function"}, {"Method(GET|POST...)", "指定的url(可选参数)", "HandlerFunc function"}]
func (service *HTTPService) BulkRouters(url string, toLowerCase bool, routers [][]interface{}) error {
	for _, val := range routers {
		valLen := len(val)
		// 执行函数取最后一个参数
		hfc := val[valLen-1].(func(http.ResponseWriter, *http.Request))
		// 最后一段路径
		var lastPath string
		if valLen < 3 {
			// 没有特殊指定, 则取函数名字
			if lastPath = service.getFunctionName(hfc); toLowerCase {
				lastPath = strings.ToLower(lastPath)
			}
		} else {
			// 特殊路径指定在第二个上
			lastPath = val[1].(string)
		}
		url := url + "/" + lastPath
		if err := service.AddHandler(val[0].(string), url, func(w http.ResponseWriter, r *http.Request) {
			hfc(w, r)
		}); nil != err {
			return err
		}
	}
	return nil
}

// AsRouter 批量注册路由, 指定一个前缀url
func (service *HTTPService) AsRouter(url string, router AsRouter) error {
	routerConfig := router.AsRouter()
	return service.BulkRouters(url, routerConfig.ToLowerCase, routerConfig.HandlerFunc)
}

// getFunctionName getFunctionName
func (service *HTTPService) getFunctionName(fun interface{}) string {
	fn := reflectutil.GetFunctionName(fun, '.')
	fm := strings.Index(fn, "-")
	if fm > -1 {
		fn = fn[:fm]
	}
	return fn
}

// Filter Filter
func (service *HTTPService) Filter(url string, fun FilterFunc) error {
	return service.AddURLFilter(url, fun)
}

// SetStaticDIR SetStaticDIR
func (service *HTTPService) SetStaticDIR(path, dir string, fun FilterFunc) (err error) {
	if !filepath.IsAbs(dir) {
		fp, err := filepath.Abs(dir)
		if nil != err {
			return err
		}
		dir = fp
	}
	handler := http.StripPrefix(path, http.FileServer(http.Dir(dir))).ServeHTTP
	if path == "/" {
		service.SetDefaultHandler(handler)
	} else {
		if err = service.AddHandler("ANY", path, handler); nil != err {
			return
		}
		if err = service.AddHandler("ANY", path+"/:**", handler); nil != err {
			return
		}
	}
	if nil != fun {
		if err = service.AddURLFilter(path, fun); nil != err {
			return
		}
		err = service.AddURLFilter(path+"/:**", fun)
	}
	return err
}

// SetStaticFile SetStaticFile
func (service *HTTPService) SetStaticFile(path, file string, fun FilterFunc) error {
	if !filepath.IsAbs(file) {
		fp, err := filepath.Abs(file)
		if nil != err {
			return err
		}
		file = fp
	}
	if err := service.AddHandler("ANY", path, func(w http.ResponseWriter, r *http.Request) {
		WirteFile(w, r, file)
	}); nil != err {
		return err
	}
	if nil != fun {
		return service.AddURLFilter(path, fun)
	}
	return nil
}

// StartHTTP 启动服务, 默认实现, 根据需要可重写
func (service *HTTPService) StartHTTP(conf StartHTTPConf) error {
	s := conf.Server
	if nil == s {
		s = &http.Server{
			WriteTimeout:   60 * time.Second,
			ReadTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}
	if nil == s.ErrorLog {
		s.ErrorLog = logs.ErrorLogger()
	}
	s.Handler = service.GetRouter()
	if len(conf.ListenAddr) > 0 {
		s.Addr = conf.ListenAddr
	}
	var err error
	if len(conf.CertFile) == 0 || len(conf.KeyFile) == 0 {
		logs.Infoln("Server(HTTP) listened in: " + s.Addr)
		err = s.ListenAndServe()
	} else {
		logs.Infoln("Server(HTTPS) listened in: " + s.Addr)
		err = s.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
	}
	return err
}
