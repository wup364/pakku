// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 模块加载器
// 依赖包: utils.utypes.Object utils.strutil.strutil

package ipakku

import (
	"reflect"

	"github.com/wup364/pakku/utils/utypes"
)

const (
	// CONST_APPNAME 示例名字
	CONST_APPNAME = "app"
	// PARAMKEY_APPNAME 实例名字KEY
	PARAMKEY_APPNAME = "app.name"
	// ErrModuleNotFoundStr 模块未找到
	ErrModuleNotFoundStr = "the module was not found, model: %s"
)

// Opts 模块配置项
type Opts struct {
	Name        string            // [可选] 模块ID, 不填则为结构体名称
	Version     float64           // [必填] 模块版本
	Description string            // [可选] 模块描述
	OnReady     func(mctx Loader) // [可选] 每次加载模块开始之前执行
	OnSetup     func()            // [可选] 模块安装, 一个模块只初始化一次
	OnUpdate    func(cv float64)  // [可选] 模块升级, 一个版本执行一次
	OnInit      func()            // [可选] 每次模块安装、升级后执行一次
}

// ModuleEvent 模块生命周期事件
type ModuleEvent string

var ModuleEventOnReady ModuleEvent = "OnReady"
var ModuleEventOnSetup ModuleEvent = "OnSetup"
var ModuleEventOnUpdate ModuleEvent = "OnUpdate"
var ModuleEventOnInit ModuleEvent = "OnInit"
var ModuleEventOnLoaded ModuleEvent = "OnLoaded"

// OnModuleEvent 模块生命周期事件回调函数
type OnModuleEvent func(module interface{}, loader Loader)

// ModuleInfo 用于记录模块信息
type ModuleInfo interface {
	Init(appName string) error
	GetValue(key string) string
	SetValue(key string, value string) error
}

// Module 实现这个接口可被加载器识别, 用于初始化和模块自动注入功能
type Module interface {
	AsModule() Opts
}

// Loader 模块加载器, 实例化后可实现统一管理模板
type Loader interface {

	// GetInstanceID 获取实例的ID
	GetInstanceID() string

	// Load 装载&初始化模块, 初始化顺序: doReady -> doSetup -> doCheckVersion -> doInit -> doEnd
	Load(mt Module)

	// Loads 装载&初始化模块, 初始化顺序: doReady -> doSetup -> doCheckVersion -> doInit -> doEnd
	Loads(mts ...Module)

	// GetParam 获取变量, 模板加载器实例上的变量
	GetParam(key string) utypes.Object

	// SetParam 设置变量, 保存在模板加载器实例内部
	SetParam(key string, val interface{})

	// GetModuleVersion 获取模块版本号
	GetModuleVersion(name string) string

	// OnModuleEvent 监听模块生命周期事件
	OnModuleEvent(name string, event ModuleEvent, val OnModuleEvent)

	// SetModuleInfoHandler 设置模块信息记录器
	SetModuleInfoHandler(moduleInfo ModuleInfo)

	// GetModuleByName 根据模块Name获取模块指针记录, 可以获取一个已经实例化的模块
	GetModuleByName(name string, val interface{}) error

	// GetModules 获取模块, 模块名字和接口名字一样才能正常获得
	GetModules(val ...interface{}) error

	// Invoke 模块调用, 返回 []reflect.Value, 返回值暂时无法处理
	Invoke(name string, method string, params ...interface{}) ([]reflect.Value, error)

	// AutoWired 自动注入依赖对象
	AutoWired(structobj interface{}) error
}
