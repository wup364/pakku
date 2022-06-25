// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of jsoncfg source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of jsoncfg software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and jsoncfg permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 入口包

package pakku

import (
	"net/http"
	"os"
	"testing"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
)

// TestNewApplication 使用现有的模块, 创建一个http服务
func TestNewApplication(t *testing.T) {
	// 实例化一个application, 启用核心模块和网络服务模板并把日志级别设置为DEBUG
	app := NewApplication("app-test").EnableCoreModule().EnableNetModule().SetLoggerLevel(logs.DEBUG).BootStart()
	// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
	var service ipakku.AppService
	if err := app.GetModuleByName("AppService", &service); nil != err {
		logs.Panicln(err)
	}
	// 手工注册一个请求路径(可使用Controller接口批量注册)
	service.SetStaticDIR("/", os.TempDir(), func(rw http.ResponseWriter, r *http.Request) bool {
		return true
	})
	service.Get("/hello", func(rw http.ResponseWriter, _ *http.Request) {
		rw.Write([]byte("hello!"))
	})
	// 启动服务
	service.StartHTTP(ipakku.HTTPServiceConfig{ListenAddr: "127.0.0.1:8080"})
	// service.StartHTTP(ipakku.HTTPServiceConfig{
	// 	KeyFile:    "./.conf/key.pem",
	// 	CertFile:   "./.conf/cert.pem",
	// 	ListenAddr: "127.0.0.1:8080",
	// })

}

// TestNewApplication1 加载自定义的模块, 创建一个http服务
func TestNewApplication1(t *testing.T) {
	// 在 TestNewApplication 示例的基础上, 新加载了一个test4Controller模块
	app := NewApplication("app-test").EnableCoreModule().EnableNetModule().AddModules(new(test4Controller)).SetLoggerLevel(logs.DEBUG).BootStart()
	// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
	var service ipakku.AppService
	if err := app.GetModuleByName("AppService", &service); nil != err {
		logs.Panicln(err)
	}
	// 启动服务
	service.StartHTTP(ipakku.HTTPServiceConfig{ListenAddr: "127.0.0.1:8080"})
}

// test4Controller 示例模块, 同时实现了Module和Controller接口
type test4Controller struct {
	//  自动注入AppService接口
	svr ipakku.AppService `@autowired:"AppService"`
}

// AsModule 作为一个模块加载
func (t *test4Controller) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:    "Test4Controller",
		Version: 1.0,
		OnInit: func() {
			if err := t.svr.AsController(t); nil != err {
				logs.Panicln(err)
			}
		},
	}
}

// AsController 实现 AsController 接口
func (ctl *test4Controller) AsController() ipakku.ControllerConfig {
	return ipakku.ControllerConfig{
		RequestMapping: "/sayhello/v1",
		RouterConfig: ipakku.RouterConfig{
			ToLowerCase: true,
			HandlerFunc: [][]interface{}{
				{http.MethodGet, "/hello", ctl.Hello},
			},
		},
	}
}

// Hello Hello
func (t *test4Controller) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Test4Controller -> Hello"))
}
