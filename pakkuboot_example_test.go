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
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/strutil"
	"github.com/wup364/pakku/utils/utypes"
)

// TestBasicNetService 使用现有的模块, 创建一个http服务
func TestBasicNetService(t *testing.T) {
	builder := NewApplication("app-example-basicnetservice")    // 实例构建器
	builder.PakkuModules().EnableAppConfig().EnableAppService() // 默认模块启用: 配置模块、网络服务模块
	builder.PakkuConfigure().SetLoggerLevel(logs.DEBUG)         // 日志级别设置为DEBUG
	app := builder.BootStart()                                  // 启动实例

	// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
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

// TestCustomModulesAndController 加载自定义的模块, 创建一个http服务
func TestCustomModulesAndController(t *testing.T) {
	// 实例化一个application, 启用配置模块、网络服务模块并加载了一个自定义模块, 把日志级别设置为DEBUG
	builder := NewApplication("app-example-custommodulesandcontroller")
	builder.PakkuModules().EnableAppConfig().EnableAppService() // 默认模块启用: 配置模块、网络服务模块
	builder.CustomModules().AddModules(new(exampleModule))      // 自定义模块加载
	builder.PakkuConfigure().SetLoggerLevel(logs.DEBUG)         // 日志级别设置为DEBUG
	app := builder.BootStart()                                  // 启动实例

	// 获取内部的一个模块, 这里使用 AppService 用于开启一个服务
	service := app.PakkuModules().GetAppService()

	// 注册一个controller
	checkError(service.AsController(new(exampleController)))

	// 启动服务
	service.StartHTTP(ipakku.HTTPServiceConfig{ListenAddr: "127.0.0.1:8080"})
}

// exampleIterface 示例模块接口
type exampleIterface interface {
	SayHello() string
}

// exampleConfigBean 测试模块配置结构体
type exampleConfigBean struct {
	exampleConfigBean1
	value10 struct {
		value2 string `@value:"value-str"`
	}
	confOther exampleConfigBean1 `@autoConfig:"other"`
	value9    utypes.Object      `@value:""`
	value8    utypes.Object      `@value:""`
	value7    float64            `@value:""`
	value6    float64            `@value:":-1"`
	value5    float64            `@value:"value-float:-1"`
	value4    []int64            `@value:"value-nums"`
	value3    []string           `@value:"value-strs"`
	value2    string             `@value:"value-str"`
	value1    int64              `@value:"value-int:-1"`
}
type exampleConfigBean1 struct {
	value2 string `@value:"value-str"`
}

// exampleModule 示例模块, 实现了Module接口
type exampleModule struct {
	config        ipakku.AppConfig  `@autowired:""`
	exampleConfig exampleConfigBean `@autoConfig:"test"`
}

// AsModule 作为一个模块加载
func (t *exampleModule) AsModule() ipakku.Opts {
	return ipakku.Opts{
		// Name:    "exampleModule",
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
func (t *exampleModule) SayHello() string {
	return fmt.Sprintf("exampleController -> Hello, conf=%v, r=%s", t.exampleConfig, strutil.GetRandom(6))
}

// printConfigs 打印输出配置信息
func (t *exampleModule) printConfigs() {
	logs.Infof("confOther: %v", t.exampleConfig.confOther)
	logs.Infof("exampleConfigBean1.value2: %v", t.exampleConfig.exampleConfigBean1.value2)
	logs.Infof("value10: %v", t.exampleConfig.value10)
	logs.Infof("value9: %v", t.exampleConfig.value9)
	logs.Infof("value8: %v", t.exampleConfig.value8)
	logs.Infof("value7: %v", t.exampleConfig.value7)
	logs.Infof("value6: %v", t.exampleConfig.value6)
	logs.Infof("value5: %v", t.exampleConfig.value5)
	logs.Infof("value4: %v", t.exampleConfig.value4)
	logs.Infof("value3: %v", t.exampleConfig.value3)
	logs.Infof("value2: %v", t.exampleConfig.value2)
	logs.Infof("value1: %v", t.exampleConfig.value1)
}

// updateAndSaveConfigs 设置一些测试数据, 下次启动或重新读取配置时则会读取出来
func (t *exampleModule) updateAndSaveConfigs() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	checkError(t.config.SetConfig("test.value-str", strutil.GetRandom(3)))
	checkError(t.config.SetConfig("test.value-int", r.Int63()))
	checkError(t.config.SetConfig("test.value-float", -r.Float64()))
	checkError(t.config.SetConfig("test.value-strs", []string{"str1", strutil.GetRandom(3)}))
	checkError(t.config.SetConfig("test.value-nums", []int64{r.Int63(), -r.Int63()}))
	checkError(t.config.SetConfig("test.value8", map[string][]int{strutil.GetRandom(3): {1}}))
	checkError(t.config.SetConfig("test.other.value-str", strutil.GetRandom(3)))
}

// reloadConfigs 重新手动将配置读取处理
func (t *exampleModule) reloadConfigs() {
	checkError(t.config.ScanAndAutoValue("test", &t.exampleConfig))
}

// exampleController 示例控制器, 实现了Controller接口
type exampleController struct {
	test exampleIterface `@autowired:"exampleModule"`
}

// AsController 实现 AsController 接口
func (ctl *exampleController) AsController() ipakku.ControllerConfig {
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
func (t *exampleController) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(t.test.SayHello()))
}
