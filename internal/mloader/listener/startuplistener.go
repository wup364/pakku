// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 模块加载监听 - 默认监听事件, 主要负责模块的初始化

package listener

import (
	"github.com/wup364/pakku/internal/mloader/mutils"
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"
)

// StartupListener 默认监听处理
type StartupListener struct{}

// Bind 绑定事件监听
func (evt *StartupListener) Bind(m ipakku.Modules) {
	m.OnModuleEvent("*", ipakku.ModuleEventOnReady, evt.doReady)
}

// doReady 模块准备
func (evt *StartupListener) doReady(m any, app ipakku.Application) {
	if err := mutils.AutoWired(m, app); nil != err {
		logs.Panic(err)
	}
}
