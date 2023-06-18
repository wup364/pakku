// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 缓存工具

package appcache

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"

	// 注册
	_ "github.com/wup364/pakku/modules/appcache/localcache"
)

// AppCache 配置模块
type AppCache struct {
	appname string
	cache   ipakku.ICache
	conf    ipakku.AppConfig `@autowired:"AppConfig"`
}

// AsModule 作为一个模块加载
func (cache *AppCache) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "AppCache",
		Version:     1.0,
		Description: "AppCache module",
		OnReady: func(mctx ipakku.Loader) {
			// 获取配置的适配器, 默认本地
			if err := ipakku.Override.AutowireInterfaceImpl(mctx, &cache.cache, "local"); nil != err {
				logs.Panicln(err)
			}
			cache.appname = mctx.GetParam(ipakku.PARAMKEY_APPNAME).ToString("app")
		},
		OnInit: func() {
			// 初始化配置
			cache.cache.Init(cache.conf, cache.appname)
		},
	}
}

// RegLib lib为库名, second:过期时间-1为不过期
func (cache *AppCache) RegLib(clib string, second int64) error {
	return cache.cache.RegLib(clib, second)
}

// Set 向lib库中设置键为key的值tb
func (cache *AppCache) Set(clib string, key string, tb interface{}) error {
	return cache.cache.Set(clib, key, tb)
}

// Get 读取缓存信息
func (cache *AppCache) Get(clib string, key string, val interface{}) error {
	return cache.cache.Get(clib, key, val)
}

// DEL 删除缓存信息
func (cache *AppCache) Del(clib string, key string) error {
	return cache.cache.Del(clib, key)
}

// Keys 获取库的所有key
func (cache *AppCache) Keys(clib string) []string {
	return cache.cache.Keys(clib)
}

// Clear 清空库内容
func (cache *AppCache) Clear(clib string) {
	cache.cache.Clear(clib)
}
