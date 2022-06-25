// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 缓存工具

package localcache

import (
	"sync"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/utypes"
)

func init() {
	ipakku.Override.RegisterInterfaceImpl(new(CacheManager), "ICache", "local")
}

// StructClone 本机cache接口, 用于交换内存数据.
type StructClone interface {
	Clone(val interface{}) error
}

// CacheManager 基于TokenManager实现的缓存管理器
// 使用前需要调用 init 方法
type CacheManager struct {
	// defaultlib string
	libexp   map[string]int64
	clibs    map[string]*TokenManager
	clibLock *sync.RWMutex
}

// Init 初始化缓存管理器, 一个对象只能初始化一次
func (cm *CacheManager) Init(config ipakku.AppConfig, appname string) {
	if nil != cm.clibs {
		return
	}
	// cm.defaultlib = "_d_"
	cm.libexp = make(map[string]int64)
	cm.clibs = make(map[string]*TokenManager)
	cm.clibLock = new(sync.RWMutex)
	// 初始默认库, 默认初始化一个名为_d_的缓存库, 该库存储的内容永不过期
	// cm.clibLock.Lock()
	// cm.clibs[cm.defaultlib] = (&TokenManager{}).Init()
	// cm.libexp[cm.defaultlib] = -1
	// cm.clibLock.Unlock()
}

// RegLib 注册缓存库
// lib为库名, second:过期时间-1为不过期
func (cm *CacheManager) RegLib(clib string, second int64) error {
	if len(clib) == 0 {
		return ipakku.ErrCacheLibNotExist
	}
	defer cm.clibLock.Unlock()
	cm.clibLock.Lock()

	if _, ok := cm.clibs[clib]; ok {
		return ipakku.ErrCacheLibIsExist
	}
	cm.clibs[clib] = (&TokenManager{}).Init()
	cm.libexp[clib] = second

	return nil
}

// Set 向lib库中设置键为key的值tb
func (cm *CacheManager) Set(clib string, key string, tb interface{}) error {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()
	if len(clib) == 0 {
		return ipakku.ErrCacheLibNotExist
	}
	tm, ok := cm.clibs[clib]
	if !ok {
		return ipakku.ErrCacheLibNotExist
	}
	lx, ok := cm.libexp[clib]
	if !ok {
		delete(cm.clibs, clib)
		return ipakku.ErrCacheLibNotExist
	}
	tm.PutTokenBody(key, tb, lx)
	return nil
}

// Get 读取缓存信息
func (cm *CacheManager) Get(clib string, key string, val interface{}) error {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	tm, ok := cm.clibs[clib]
	if !ok {
		return ipakku.ErrCacheLibNotExist
	}
	if tmp, ok := tm.GetTokenBody(key); ok && nil != val {
		if clone, ok := tmp.(StructClone); ok {
			return clone.Clone(val)
		}
		return utypes.NewObject(tmp).Scan(val)
	}
	return ipakku.ErrNoCacheHit
}

// Del 删除缓存信息
func (cm *CacheManager) Del(clib string, key string) error {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	if tm, ok := cm.clibs[clib]; !ok {
		return ipakku.ErrCacheLibNotExist
	} else {
		tm.DestroyToken(key)
	}
	return nil
}

// Keys 获取库的所有key
func (cm *CacheManager) Keys(clib string) []string {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	tm, ok := cm.clibs[clib]
	if !ok {
		return make([]string, 0)
	}
	return tm.ListTokens()
}

// Clear 清空库内容
func (cm *CacheManager) Clear(clib string) {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	if tm, ok := cm.clibs[clib]; ok {
		tm.Clear()
	}
}
