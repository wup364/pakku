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

package mloader

import (
	"testing"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
)

// DemoModule 示例模块
type DemoModule struct {
}

// AsModule 作为一个模块加载
func (t *DemoModule) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "DemoModule",
		Version:     1.0,
		Description: "示例模板",
		Updaters:    func(app ipakku.Application) ipakku.Updaters { return make([]ipakku.Updater, 0) },
		OnReady: func(app ipakku.Application) {
			logs.Infoln("on ready")
		},
		OnSetup: func() {
			logs.Infoln("on setup")
		},
		OnInit: func() {
			logs.Infoln("on init")
		},
	}
}

// Hello Hello
func (t *DemoModule) Hello() {
	logs.Infoln("DemoModule -> Hello")
}

// 在 mian 中调用
func TestLoader(t *testing.T) {
	loader := NewDefault("Test")
	// loader.SetModuleInfoRecorder(xxx)
	loader.Loads(new(DemoModule))
	loader.GetApplication().Utils().Invoke("DemoModule", "Hello")
}
