// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 模块加载器
// 依赖包: utils.utypes.Object utils.strutil.strutil

package mloader

import (
	"testing"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"
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
			logs.Info("on ready")
		},
		OnSetup: func() {
			logs.Info("on setup")
		},
		OnInit: func() {
			logs.Info("on init")
		},
	}
}

// Hello Hello
func (t *DemoModule) Hello() {
	logs.Info("DemoModule -> Hello")
}

// 在 mian 中调用
func TestLoader(t *testing.T) {
	loader := NewDefault("Test")
	// loader.SetModuleInfoRecorder(xxx)
	loader.Loads(new(DemoModule))
	loader.GetApplication().Utils().Invoke("DemoModule", "Hello")
}
