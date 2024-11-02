// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package ipakku

import (
	"net"
	"net/http"

	"github.com/wup364/pakku/utils/serviceutil"
)

const (
	// CONFKEY_READTIMEOUTSECOND ReadTimeoutSecond
	CONFKEY_READTIMEOUTSECOND = "service.ReadTimeoutSecond"
	// CONFKEY_WRITETIMEOUTSECOND WriteTimeoutSecond
	CONFKEY_WRITETIMEOUTSECOND = "service.WriteTimeoutSecond"
	// CONFKEY_MAXHEADERBYTES MaxHeaderBytes
	CONFKEY_MAXHEADERBYTES = "service.MaxHeaderBytes"
)

// HandlerFunc 定义请求处理器
type HandlerFunc func(http.ResponseWriter, *http.Request)

// FilterFunc http请求过滤器, 返回bool, true: 继续, false: 停止
type FilterFunc func(http.ResponseWriter, *http.Request) bool

// Filter4Passed 空过滤器(通过的): 没有任何处理逻辑的过滤器
var Filter4Passed FilterFunc = func(http.ResponseWriter, *http.Request) bool { return true }

// RouterConfig 批量注册服务路径配置对象;
// ToLowerCase bool   是否需要url转小写, 在未指定url(使用函数名字作为url的一部分)的情况下生效;
// HandlerFunc [][]any 需要注册的函数 [{"Method(GET|POST...)", "HandlerFunc function"}, {"Method(GET|POST...)", "指定的url(可选参数)", "HandlerFunc function"}]
type RouterConfig serviceutil.RouterConfig

// FilterConfigItem 过滤器配置对象
type FilterConfigItem struct {
	Path string     // 指定的url, 自动添加RequestMapping前缀
	Func FilterFunc // http请求过滤器, 返回bool, true: 继续, false: 停止
}

// ControllerConfig  注册对象为Controller配置对象;
// RequestMapping 请求路径, 也可以是版本号(v1|v2...)作为路径的一部分;
// RouterConfig 批量注册服务路径配置对象
type ControllerConfig struct {
	RequestMapping string             // 请求路径, 也可以是版本号(v1|v2...)作为路径的一部分;
	RouterConfig                      // 批量注册服务路径配置对象
	FilterConfig   []FilterConfigItem // 过滤器配置对象, 自动添加前缀路径(RequestMapping值)
}

// Router 批量注册服务路径
type Router interface {
	AsRouter() serviceutil.RouterConfig
}

// Controller 注册对象为Controller
type Controller interface {
	AsController() ControllerConfig
}

// HTTP 方法 GET POST .....
type Method string

// HTTPService 服务
type HTTPService interface {
	// Get Get
	Get(url string, fun HandlerFunc) error

	// Post Post
	Post(url string, fun HandlerFunc) error

	// Put Put
	Put(url string, fun HandlerFunc) error

	// Patch Patch
	Patch(url string, fun HandlerFunc) error

	// Head Head
	Head(url string, fun HandlerFunc) error

	// Options Options
	Options(url string, fun HandlerFunc) error

	// Delete Delete
	Delete(url string, fun HandlerFunc) error

	// Any Any
	Any(url string, fun HandlerFunc) error

	// AsRouter 批量注册路由, 可以再指定一个前缀url
	AsRouter(url string, router Router) error

	// AsController 批量注册路由, 使用RequestMapping字段作为前缀url
	AsController(router Controller) error

	// Filter 注册请求过滤器
	Filter(url string, fun FilterFunc) error

	// SetStaticDIR SetStaticDIR
	SetStaticDIR(path, dir string, fun FilterFunc) error

	// SetStaticFile SetStaticFile
	SetStaticFile(path, file string, fun FilterFunc) error
}

// RPCService 服务
type RPCService interface {
	RegisteRPC(rcvr any) error
}

// HTTPServiceConfig 启动配置
type HTTPServiceConfig struct {
	Debug      bool
	CertFile   string
	KeyFile    string
	ListenAddr string
	Server     *http.Server
}

// RPCServiceConfig 启动配置
type RPCServiceConfig struct {
	Debug      bool
	Network    string
	ListenAddr string
	Listener   net.Listener
}

// AppService web服务即接口
type AppService interface {
	HTTPService
	RPCService
	StartHTTP(serviceCfg HTTPServiceConfig)
	StartRPC(serviceCfg RPCServiceConfig)
}
