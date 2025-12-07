// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 模块加载监听
package listener

import (
	"github.com/wup364/pakku/ipakku"
)

// defaultListeners 默认已注册的监听(在初始化时注册, 比所有模块都要早执行), 在启动时按照顺序加载
var defaultListeners = []MListener{
	new(StartupListener),
}

// MListener 模块加载器监听
type MListener interface {
	Bind(m ipakku.Modules)
}

// GetDefaultListeners 获取默认已注册的监听(在初始化时注册, 比所有模块都要早执行)
func GetDefaultListeners() []MListener {
	return defaultListeners
}
