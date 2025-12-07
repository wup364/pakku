// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 模块加载器-实例化

package mloader

import (
	"github.com/wup364/pakku/internal/mloader/listener"
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/strutil"
	"github.com/wup364/pakku/pkg/utypes"
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
