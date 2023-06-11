// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of jsoncfg source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of jsoncfg software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and jsoncfg permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// pakku演示用例

package pakku

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/strutil"
)

// TestBasicNetService 使用现有的模块, 创建一个http服务
func TestBasicNetService(t *testing.T) {
	// 实例化一个application, 启用核心模块和网络服务模板并把日志级别设置为DEBUG
	app := NewApplication("app-demo").EnableCoreModule().EnableNetModule().SetLoggerLevel(logs.DEBUG).BootStart()

	// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
	var service ipakku.AppService
	if err := app.GetModules(&service); nil != err {
		logs.Panicln(err)
	}

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

// TestCustomModulesAndController 加载自定义的模块, 创建一个http服务
func TestCustomModulesAndController(t *testing.T) {
	// 实例化一个application, 启用核心模块和网络服务模板并加载了一个自定义模块, 最后把日志级别设置为DEBUG
	app := NewApplication("app-demo").EnableCoreModule().EnableNetModule().AddModules(new(demoModule)).SetLoggerLevel(logs.DEBUG).BootStart()

	// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
	var service ipakku.AppService
	if err := app.GetModules(&service); nil != err {
		logs.Panicln(err)
	}

	// 注册一个controller
	checkError(service.AsController(new(demoController)))

	// 启动服务
	service.StartHTTP(ipakku.HTTPServiceConfig{ListenAddr: "127.0.0.1:8080"})
}

// demoIterface 示例模块接口
type demoIterface interface {
	SayHello() string
}

// demoConfigBean 测试模块配置结构体
type demoConfigBean struct {
	value7 float64  `@value:""`
	value6 float64  `@value:":-1"`
	value5 float64  `@value:"value-float:-1"`
	value4 []int64  `@value:"value-nums"`
	value3 []string `@value:"value-strs"`
	value2 string   `@value:"value-str"`
	value1 int64    `@value:"value-int:-1"`
}

// demoModule 示例模块, 实现了Module接口
type demoModule struct {
	config     ipakku.AppConfig `@autowired:"AppConfig"`
	demoConfig *demoConfigBean  `@autoConfig:"test"`
}

// AsModule 作为一个模块加载
func (t *demoModule) AsModule() ipakku.Opts {
	return ipakku.Opts{
		// Name:    "demoModule",
		Version: 1.0,
		OnInit: func() {
			t.printConfigs()
			t.updateAndSaveConfigs()
			t.reloadConfigs()
			t.printConfigs()
		},
	}
}

// SayHello SayHello
func (t *demoModule) SayHello() string {
	return fmt.Sprintf("demoController -> Hello, conf=%v, r=%s", t.demoConfig, strutil.GetRandom(6))
}

// printConfigs 打印输出配置信息
func (t *demoModule) printConfigs() {
	logs.Infof("value7: %v", t.demoConfig.value7)
	logs.Infof("value6: %v", t.demoConfig.value6)
	logs.Infof("value5: %v", t.demoConfig.value5)
	logs.Infof("value4: %v", t.demoConfig.value4)
	logs.Infof("value3: %v", t.demoConfig.value3)
	logs.Infof("value2: %v", t.demoConfig.value2)
	logs.Infof("value1: %v", t.demoConfig.value1)
}

// updateAndSaveConfigs 设置一些测试数据, 下次启动或重新读取配置时则会读取出来
func (t *demoModule) updateAndSaveConfigs() {
	checkError(t.config.SetConfig("test.value-str", strutil.GetRandom(3)))
	checkError(t.config.SetConfig("test.value-int", "-1024"))
	checkError(t.config.SetConfig("test.value-float", "-1024.1024"))
	checkError(t.config.SetConfig("test.value-strs", []string{"str1", strutil.GetRandom(3)}))
	checkError(t.config.SetConfig("test.value-nums", []int64{1204, -1024}))
}

// reloadConfigs 重新手动将配置读取处理
func (t *demoModule) reloadConfigs() {
	checkError(t.config.ScanAndAutoValue("test", t.demoConfig))
}

// demoController 示例控制器, 实现了Controller接口
type demoController struct {
	test demoIterface `@autowired:"demoModule"`
}

// AsController 实现 AsController 接口
func (ctl *demoController) AsController() ipakku.ControllerConfig {
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
func (t *demoController) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(t.test.SayHello()))
}
