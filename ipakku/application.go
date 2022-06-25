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
	"io"
	"reflect"

	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/utypes"
)

// ApplicationBoot 应用初始化引导
type ApplicationBoot interface {

	// BootStart 加载&启动程序
	BootStart() Application

	// GetApplication 获取Application
	GetApplication() Application

	// AddModule 添加模块
	AddModule(mt Module) ApplicationBoot

	// AddModules 添加模块
	AddModules(mts ...Module) ApplicationBoot

	// EnableCoreModule 启用默认的核心模块
	EnableCoreModule() ApplicationBoot

	// EnableNetModule 启用默认的网络服务模块
	EnableNetModule() ApplicationBoot

	// SetLoggerOutput 设置日志输出方式
	SetLoggerOutput(w io.Writer) ApplicationBoot

	// SetLoggerLevel 设置日志输出级别 NONE DEBUG INFO ERROR
	SetLoggerLevel(lv logs.LoggerLeve) ApplicationBoot

	// DisableBanner 禁止Banner输出
	DisableBanner() ApplicationBoot
}

// Application 应用实例, 继承 loader.Loader
type Application interface {

	// GetInstanceID 获取实例的ID
	GetInstanceID() string

	// LoadModule 装载&初始化模块 - DO Setup -> Check Ver -> Do Init
	LoadModule(mt Module) Application

	// LoadModules 装载&初始化模块 - DO Setup -> Check Ver -> Do Init
	LoadModules(mts ...Module) Application

	// GetParam 获取变量, 模板加载器实例上的变量
	GetParam(key string) utypes.Object

	// SetParam 设置变量, 保存在模板加载器实例内部
	SetParam(key string, val interface{})

	// OnModuleEvent 监听模块生命周期事件
	OnModuleEvent(name string, event ModuleEvent, val OnModuleEvent)

	// GetModuleByName 根据模块Name获取模块指针记录, 可以获取一个已经实例化的模块
	GetModuleByName(name string, val interface{}) error

	// Invoke 模块调用, 返回 []reflect.Value, 返回值暂时无法处理
	Invoke(name string, method string, params ...interface{}) ([]reflect.Value, error)

	// AutoWired 自动注入依赖对象
	AutoWired(structobj interface{}) error
}
