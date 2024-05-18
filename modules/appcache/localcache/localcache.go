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
	"fmt"
	"sync"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/utypes"
)

func init() {
	ipakku.PakkuConf.RegisterPakkuModuleImplement(new(CacheManager), "ICache", "local")
}

// StructClone 本机cache接口, 用于交换内存数据.
type StructClone interface {
	Clone(val interface{}) error
}

// StructValue 值
type StructValue struct {
	Value any
}

// Clone 本机cache接口, 用于交换内存数据.
func (sv *StructValue) Clone(val interface{}) error {
	if uat, ok := val.(*StructValue); ok {
		uat.Value = sv.Value
		return nil
	}
	return fmt.Errorf("can't support clone %T ", val)
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

// Exists 返回key是否存在
func (cm *CacheManager) Exists(clib string, key string) (res bool, err error) {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	if tm, ok := cm.clibs[clib]; !ok {
		return false, ipakku.ErrCacheLibNotExist
	} else {
		_, res = tm.GetTokenBody(key)
	}
	return
}

// Set 向lib库中设置键为key的值
// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
func (cm *CacheManager) Set(clib string, key string, args ...any) error {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	if len(clib) == 0 {
		return ipakku.ErrCacheLibNotExist
	}

	if tm, ok := cm.clibs[clib]; !ok {
		return ipakku.ErrCacheLibNotExist
	} else if lx, err := cm.getExpSecond(clib, args...); nil != err {
		return err
	} else {
		tm.PutTokenBody(key, args[0], lx)
		return nil
	}
}

// SetNX 向lib库中设置键为key的值, 当key不存在时设置成功, 并返回true
// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
func (cm *CacheManager) SetNX(clib string, key string, args ...any) (bool, error) {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	if len(clib) == 0 {
		return false, ipakku.ErrCacheLibNotExist
	}

	if tm, ok := cm.clibs[clib]; !ok {
		return false, ipakku.ErrCacheLibNotExist
	} else if lx, err := cm.getExpSecond(clib, args...); nil != err {
		return false, err
	} else {
		return tm.PutTokenBodyNX(key, args[0], lx), nil
	}
}

// Incrby 指定key以increment的值累加, 返回累加后的值
// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
func (cm *CacheManager) Incrby(clib string, key string, args ...any) (int64, error) {
	defer cm.clibLock.RUnlock()
	cm.clibLock.RLock()

	if len(clib) == 0 {
		return -1, ipakku.ErrCacheLibNotExist
	}

	var ok bool
	var tm *TokenManager
	if tm, ok = cm.clibs[clib]; !ok {
		return -1, ipakku.ErrCacheLibNotExist
	}

	// 查询已存在
	val := args[0].(int64)
	if oldval, ok := tm.GetTokenBody(key); ok {
		tmp := oldval.(*StructValue)
		tmp.Value = tmp.Value.(int64) + val
		return tmp.Value.(int64), nil
	}

	// 第一次插入
	if lx, err := cm.getExpSecond(clib, args...); nil != err {
		return -1, err
	} else {
		val := args[0].(int64)
		tm.PutTokenBody(key, &StructValue{val}, lx)
		return val, nil
	}
}

// getExpSecond 获取过期时间
func (cm *CacheManager) getExpSecond(clib string, args ...any) (int64, error) {
	if len(args) > 1 {
		if val, ok := args[1].(int64); ok {
			return val, nil
		} else if val, ok := args[1].(int); ok {
			return int64(val), nil
		} else {
			return -1, ipakku.ErrCacheArgsTypeError
		}
	} else if lx, ok := cm.libexp[clib]; !ok {
		delete(cm.clibs, clib)
		return -1, ipakku.ErrCacheLibNotExist
	} else {
		return lx, nil
	}
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
