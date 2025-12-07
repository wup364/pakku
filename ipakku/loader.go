// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 模块加载器
// 依赖包: utils.utypes.Object utils.strutil.strutil

package ipakku

import (
	"reflect"

	"github.com/wup364/pakku/pkg/utypes"
)

const (
	// DEFT_VAL_APPNAME 默认实例名字
	DEFT_VAL_APPNAME = "app"
	// PARAMS_KEY_APPNAME 实例名字KEY
	PARAMS_KEY_APPNAME = "app.name"
	// ERR_MSG_MODULE_NOT_FOUND 模块未找到
	ERR_MSG_MODULE_NOT_FOUND = "the module was not found, model: %s"
)

// Updater 模块版本升级执行器
type Updater interface {
	// Version 要升级到的版本号
	Version() float64
	// Execute 执行升级
	Execute(app Application) error
}

// Updaters 升级器
type Updaters []Updater

// 实现sort.Interface接口取元素数量方法
func (sort Updaters) Len() int {
	return len(sort)
}

// 实现sort.Interface接口比较元素方法
func (sort Updaters) Less(i, j int) bool {
	return sort[i].Version() < sort[j].Version()
}

// 实现sort.Interface接口交换元素方法
func (sort Updaters) Swap(i, j int) {
	sort[i], sort[j] = sort[j], sort[i]
}

// Opts 模块配置项
type Opts struct {
	Name        string                         // [可选] 模块ID, 不填则为结构体名称
	Version     float64                        // [必填] 模块版本
	Description string                         // [可选] 模块描述
	Updaters    func(app Application) Updaters // [可选] 模块升级执行器, 一个版本执行一次
	OnReady     func(app Application)          // [可选] 每次加载模块开始之前执行
	OnSetup     func()                         // [可选] 模块安装, 一个模块只初始化一次
	OnInit      func()                         // [可选] 每次模块安装、升级后执行一次
}

// ModuleEvent 模块生命周期事件
type ModuleEvent string

var ModuleEventOnReady ModuleEvent = "OnReady"
var ModuleEventOnSetup ModuleEvent = "OnSetup"
var ModuleEventOnUpdate ModuleEvent = "OnUpdate"
var ModuleEventOnInit ModuleEvent = "OnInit"
var ModuleEventOnLoaded ModuleEvent = "OnLoaded"

var ModuleEventOnSetupSucced ModuleEvent = "OnSetupSucced"
var ModuleEventOnUpdateSucced ModuleEvent = "OnUpdateSucced"

// OnModuleEvent 模块生命周期事件回调函数
type OnModuleEvent func(module any, app Application)

// ModuleInfoRecorder 用于记录模块信息
type ModuleInfoRecorder interface {
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

	// Loads 装载&初始化模块(自动分析模块依赖顺序), 初始化顺序: doReady -> doSetup -> doCheckVersion -> doInit -> doEnd
	Loads(mts ...Module)

	// SetModuleInfoRecorder 设置模块信息记录器
	SetModuleInfoRecorder(moduleInfo ModuleInfoRecorder)

	// Application 获取当前实例
	GetApplication() Application

	Params  // Params 保存实例中的键值对数据
	Modules // Modules 模块操作
}

// Application 当前运行中的实例
type Application interface {
	// GetInstanceID 获取实例的ID
	GetInstanceID() string

	// Params 保存实例中的键值对数据
	Params() Params

	// Modules 模块操作
	Modules() Modules

	// Utils 工具
	Utils() Utils
}

// Params 保存实例中的键值对数据
type Params interface {
	ParamGetter
	ParamSetter
}

// ParamGetter 只读 - 保存实例中的键值对数据
type ParamGetter interface {
	// GetParam 获取变量, 当前实例上的变量
	GetParam(key string) utypes.Object
}

// ParamSetter 只写 - 保存实例中的键值对数据
type ParamSetter interface {
	// SetParam 设置变量, 保存在当前实例内部
	SetParam(key string, val any)
}

// Modules 模块操作
type Modules interface {
	// GetModuleByName 根据模块Name获取模块指针记录, 可以获取一个已经实例化的模块
	GetModuleByName(name string, val any) error

	// GetModules 获取模块, 模块名字和接口名字一样才能正常获得
	GetModules(val ...any) error

	// GetModuleVersion 获取模块版本号
	GetModuleVersion(name string) string

	// OnModuleEvent 监听模块生命周期事件
	OnModuleEvent(name string, event ModuleEvent, val OnModuleEvent)
}

// Utils 工具
type Utils interface {
	// AutoWired 自动注入依赖对象
	AutoWired(structobj any) error

	// Invoke 模块调用, 返回 []reflect.Value, 返回值暂时无法处理
	Invoke(name string, method string, params ...any) ([]reflect.Value, error)
}
