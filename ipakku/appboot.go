// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package ipakku

import (
	"io"

	"github.com/wup364/pakku/pkg/logs"
)

// ApplicationBootBuilder 应用初始化引导
type ApplicationBootBuilder interface {

	// PakkuConfigure 应用配置操作
	PakkuConfigure() PakkuConfigure

	// PakkuModules 默认模块启用操作
	PakkuModules() PakkuModuleBuilder

	// CustomModules 自定义模块操作
	CustomModules() CustomModuleBuilder

	// ModuleEvents 模块事件监听器
	ModuleEvents() ModuleEventBuilder

	// Application 获取Application实例
	Application() PakkuApplication

	// BootStart 加载&启动程序
	BootStart() PakkuApplication
}

// PakkuApplication bootBuild实例化后的application
type PakkuApplication interface {
	Application

	// PakkuModules 默认模块Getter
	PakkuModules() PakkuModulesGetter
}

// PakkuModule 应用配置
type PakkuConfigure interface {
	// SetLoggerOutput 设置日志输出方式
	SetLoggerOutput(w io.Writer) PakkuConfigure

	// SetLoggerLevel 设置日志输出级别 NONE DEBUG INFO ERROR
	SetLoggerLevel(lv logs.LoggerLeve) PakkuConfigure

	// DisableBanner 禁止Banner输出
	DisableBanner() PakkuConfigure

	// PakkuModules 启用默认携带的模块
	PakkuModules() PakkuModuleBuilder

	// CustomModules 自定义模块操作
	CustomModules() CustomModuleBuilder
}

// PakkuModule 默认模块启用操作
type PakkuModuleBuilder interface {

	// EnableAppConfig 启用配置模块
	EnableAppConfig() PakkuModuleBuilder

	// EnableAppCache 启用缓存模块
	EnableAppCache() PakkuModuleBuilder

	// EnableAppEvent 启用事件模块
	EnableAppEvent() PakkuModuleBuilder

	// EnableAppService 启用网络服务[WEB|RPC]模块
	EnableAppService() PakkuModuleBuilder

	// EnableAppStaticPage 启用静态页面模块, 依赖 AppService 模块
	EnableAppStaticPage() PakkuModuleBuilder

	// CustomModules 自定义模块操作
	CustomModules() CustomModuleBuilder

	// ModuleEvents 模块事件监听器
	ModuleEvents() ModuleEventBuilder

	// BootStart 加载&启动程序
	BootStart() PakkuApplication
}

// PakkuModulesGetter 获取默认携带的模块
type PakkuModulesGetter interface {

	// GetAppConfig 获得配置模块
	GetAppConfig() AppConfig

	// GetAppCache 获得缓存模块
	GetAppCache() AppCache

	// GetAppEvent 获得事件模块
	GetAppEvent() AppEvent

	// GetAppService 获得网络服务[WEB|RPC]模块
	GetAppService() AppService
}

// CustomModuleBuilder 自定义模块操作
type CustomModuleBuilder interface {

	// AddModule 添加模块
	AddModule(mt Module) CustomModuleBuilder

	// AddModules 添加模块
	AddModules(mts ...Module) CustomModuleBuilder

	// ModuleEvents 模块事件监听器
	ModuleEvents() ModuleEventBuilder

	// BootStart 加载&启动程序
	BootStart() PakkuApplication
}

// ModuleEventBuilder 模块事件监听器
type ModuleEventBuilder interface {

	// Listen 监听模块生命周期事件
	Listen(name string, event ModuleEvent, val OnModuleEvent) ModuleEventBuilder

	// BootStart 加载&启动程序
	BootStart() PakkuApplication
}
