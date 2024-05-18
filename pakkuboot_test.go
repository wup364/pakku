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
	"testing"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
)

// TestNewApplication 使用现有的模块, 创建一个http服务
func TestNewApplication(t *testing.T) {
	builder := NewApplication("app-example-basicnetservice") // 实例构建器
	app := builder.
		PakkuConfigure().SetLoggerLevel(logs.DEBUG).         // 日志级别设置为DEBUG
		PakkuModules().EnableAppConfig().EnableAppService(). // 默认模块启用: 配置模块、网络服务模块
		// CustomModules().AddModule(new(exampleModule)).    // 自定义模块加载
		// ModuleEvents().Listen("", ipakku.ModuleEventOnLoaded, func(module interface{}, app ipakku.Application) {}).
		BootStart() // 启动实例

	// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
	// var service ipakku.AppService
	// if err := app.Modules().GetModules(&service); nil != err {
	// 	logs.Panicln(err)
	// }
	service := app.PakkuModules().GetAppService()

	// 设置一个静态页面路径
	if err := service.SetStaticDIR("/", "./", nil); nil != err {
		logs.Panicln(err)
	}

	// 手工注册一个请求路径(可使用Controller接口批量注册)
	if err := service.Get("/hello", func(rw http.ResponseWriter, _ *http.Request) {
		rw.Write([]byte("hello!"))
	}); nil != err {
		logs.Panicln(err)
	}

	// 启动服务
	service.StartHTTP(ipakku.HTTPServiceConfig{ListenAddr: "127.0.0.1:8080"})
	// service.StartHTTP(ipakku.HTTPServiceConfig{
	// 	KeyFile:    "./.conf/key.pem",
	// 	CertFile:   "./.conf/cert.pem",
	// 	ListenAddr: "127.0.0.1:8080",
	// })

}

func checkError(err error) {
	if nil != err {
		panic(err)
	}
}
