// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 缓存工具

package appcache

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"

	// 注册
	_ "github.com/wup364/pakku/internal/modules/appcache/localcache"
)

// AppCache 配置模块
type AppCache struct {
	appname string
	cache   ipakku.ICache
	conf    ipakku.AppConfig `@autowired:""`
}

// AsModule 作为一个模块加载
func (cache *AppCache) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Version:     1.0,
		Description: "AppCache module",
		OnReady: func(app ipakku.Application) {
			// 获取配置的适配器, 默认本地
			if err := ipakku.PakkuConf.AutowirePakkuModuleImplement(app.Params(), &cache.cache, "local"); nil != err {
				logs.Panic(err)
			}
			cache.appname = app.Params().GetParam(ipakku.PARAMS_KEY_APPNAME).ToString(ipakku.DEFT_VAL_APPNAME)
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

// Exists 返回key是否存在
func (cache *AppCache) Exists(clib string, key string) (bool, error) {
	return cache.cache.Exists(clib, key)
}

// Set 向lib库中设置键为key的值
// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
func (cache *AppCache) Set(clib string, key string, args ...any) error {
	if len(args) < 1 {
		return ipakku.ErrCacheArgsEmpty
	}
	return cache.cache.Set(clib, key, args...)
}

// SetNX 向lib库中设置键为key的值, 当key不存在时设置成功, 并返回true
// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
func (cache *AppCache) SetNX(clib string, key string, args ...any) (bool, error) {
	if len(args) < 1 {
		return false, ipakku.ErrCacheArgsEmpty
	}
	return cache.cache.SetNX(clib, key, args...)
}

// Incrby 指定key以increment的值累加, 返回累加后的值
// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
func (cache *AppCache) Incrby(clib string, key string, args ...any) (int64, error) {
	if len(args) < 1 {
		return -1, ipakku.ErrCacheArgsEmpty
	} else if val, ok := args[0].(int); ok {
		args[0] = int64(val)
	} else if _, ok := args[0].(int64); !ok {
		return -1, ipakku.ErrCacheArgsTypeError
	}
	return cache.cache.Incrby(clib, key, args...)
}

// Get 读取缓存信息
func (cache *AppCache) Get(clib string, key string, val any) error {
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
