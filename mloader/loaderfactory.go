// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 模块加载器-实例化

package mloader

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/mloader/listener"
	"github.com/wup364/pakku/utils/strutil"
	"github.com/wup364/pakku/utils/utypes"
)

// NewDefault 实例一个带默认功能的加载器对象
func NewDefault(name string) ipakku.Loader {
	loader := New(name)
	loader.SetParam(ipakku.PARAMS_KEY_APPNAME, name)
	loader.SetModuleInfoRecorder(ipakku.PakkuConf.GetModuleInfoRecorderImplement())

	// 默认监听器列表
	if listeners := listener.GetDefaultListeners(); len(listeners) > 0 {
		for _, v := range listeners {
			v.Bind(loader)
		}
	}
	return loader
}

// NewDefault 实例一个加载器对象
func New(name string) ipakku.Loader {
	return &Loader{
		events:     utypes.NewSafeMap[string, []ipakku.OnModuleEvent](),
		modules:    utypes.NewSafeMap[string, ipakku.Module](),
		mparams:    utypes.NewSafeMap[string, any](),
		instanceID: strutil.GetUUID(),
	}
}
