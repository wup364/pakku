// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package ipakku

import "github.com/wup364/pakku/pkg/utypes"

// AppConfig app 配置模块
type AppConfig interface {

	// GetConfig 读取key的value信息, 返回 Object 对象, 里面的值可能是string或者map
	GetConfig(key string) utypes.Object

	// SetConfig 设置值
	SetConfig(key string, value any) error

	// ScanAndAutoConfig 扫描带有@autoconfig标签的字段, 并完成其配置
	ScanAndAutoConfig(ptr any) error

	// ScanAndAutoValue 扫描带有@value标签的字段, 并完成其配置
	ScanAndAutoValue(configPrefix string, ptr any) error
}

// IConfig 配置接口
type IConfig interface {

	// Init 初始化解析器
	Init(appName string) error

	// GetConfig 读取key的value信息, 返回 Object 对象, 里面的值可能是string或者map
	GetConfig(key string) (res utypes.Object)

	// SetConfig 设置值
	SetConfig(key string, value any) error
}
