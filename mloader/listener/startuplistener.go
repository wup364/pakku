// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 模块加载监听 - 默认监听事件, 主要负责模块的初始化

package listener

import (
	"strconv"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/mloader/mutils"
	"github.com/wup364/pakku/utils/logs"
)

// StartupListener 默认监听处理
type StartupListener struct{}

// Bind 绑定事件监听
func (evt *StartupListener) Bind(l ipakku.Loader) {
	l.OnModuleEvent("*", ipakku.ModuleEventOnReady, evt.doReady)
	l.OnModuleEvent("*", ipakku.ModuleEventOnSetup, evt.doSetup)
	l.OnModuleEvent("*", ipakku.ModuleEventOnUpdate, evt.doUpdate)
	l.OnModuleEvent("*", ipakku.ModuleEventOnInit, evt.doInit)
}

// doReady 模块准备
func (evt *StartupListener) doReady(m interface{}, l ipakku.Loader) {
	mt := m.(ipakku.Module)
	if err := mutils.AutoWired(mt, l); nil != err {
		logs.Panicln(err)
	}
	if nil != mt.AsModule().OnReady {
		logs.Infof("> Execute Module.OnReady \r\n")
		mt.AsModule().OnReady(l)
	}
}

// doSetup 模块安装
func (evt *StartupListener) doSetup(m interface{}, l ipakku.Loader) {
	mt := m.(ipakku.Module)
	if nil != mt.AsModule().OnSetup {
		logs.Infof("> Execute Module.OnSetup \r\n")
		mt.AsModule().OnSetup()
	}
}

// doUpdate 模块升级
func (evt *StartupListener) doUpdate(m interface{}, l ipakku.Loader) {
	mt := m.(ipakku.Module)
	if nil != mt.AsModule().OnUpdate {
		logs.Infof("> Execute Module.OnUpdate \r\n")
		if hv, err := strconv.ParseFloat(l.GetVersion(mt.AsModule().Name), 64); nil != err {
			logs.Panicln(err)
		} else {
			mt.AsModule().OnUpdate(hv)
		}
	}
}

// doInit 模块初始化
func (evt *StartupListener) doInit(m interface{}, l ipakku.Loader) {
	mt := m.(ipakku.Module)
	if nil != mt.AsModule().OnInit {
		logs.Infof("> Execute Module.OnInit \r\n")
		mt.AsModule().OnInit()
	}
}
