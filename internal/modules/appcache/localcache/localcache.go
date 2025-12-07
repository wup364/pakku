// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 缓存工具
package localcache

import (
	"sync"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/utypes"
)

func init() {
	ipakku.PakkuConf.RegisterPakkuModuleImplement(new(CacheManager), "ICache", "local")
}

// CacheValue 缓存值对象转换接口, 缓存值若实现此接口, Get时会调用
type CacheValue interface {
	LocalCacheValueScan(val any) error
}

// CacheManager 基于TokenManager实现的缓存管理器
// 使用前需要调用 init 方法
type CacheManager struct {
	libexp map[string]int64
	clibs  map[string]*TokenManager
	locker *sync.RWMutex
}

// Init 初始化缓存管理器, 一个对象只能初始化一次
func (cm *CacheManager) Init(config ipakku.AppConfig, appname string) {
	if nil != cm.clibs {
		return
	}
	cm.libexp = make(map[string]int64)
	cm.clibs = make(map[string]*TokenManager)
	cm.locker = new(sync.RWMutex)
}

// RegLib 注册缓存库
// lib为库名, second:过期时间-1为不过期
func (cm *CacheManager) RegLib(clib string, second int64) error {
	if len(clib) == 0 {
		return ipakku.ErrCacheLibNotExist
	}
	defer cm.locker.Unlock()
	cm.locker.Lock()

	if _, ok := cm.clibs[clib]; ok {
		return ipakku.ErrCacheLibIsExist
	}
	cm.clibs[clib] = (&TokenManager{}).Init()
	cm.libexp[clib] = second

	return nil
}

// Exists 返回key是否存在
func (cm *CacheManager) Exists(clib string, key string) (res bool, err error) {
	defer cm.locker.RUnlock()
	cm.locker.RLock()

	if tm, ok := cm.clibs[clib]; !ok {
		return false, ipakku.ErrCacheLibNotExist
	} else {
		_, res = tm.GetTokenBody(key)
	}
	return
}

// Get 读取缓存信息
func (cm *CacheManager) Get(clib string, key string, val any) error {
	defer cm.locker.RUnlock()
	cm.locker.RLock()

	tm, ok := cm.clibs[clib]
	if !ok {
		return ipakku.ErrCacheLibNotExist
	}
	if tmp, ok := tm.GetTokenBody(key); ok && nil != val {
		if cv, ok := tmp.(CacheValue); ok {
			return cv.LocalCacheValueScan(val)
		}
		return utypes.NewObject(tmp).Scan(val)
	}
	return ipakku.ErrNoCacheHit
}

// Keys 获取库的所有key
func (cm *CacheManager) Keys(clib string) []string {
	defer cm.locker.RUnlock()
	cm.locker.RLock()

	tm, ok := cm.clibs[clib]
	if !ok {
		return make([]string, 0)
	}
	return tm.ListTokens()
}

// Set 向lib库中设置键为key的值
// args[0] 为缓存值 args[1]如果存在, 则覆盖默认过期时间, 单位秒
func (cm *CacheManager) Set(clib string, key string, args ...any) error {
	defer cm.locker.Unlock()
	cm.locker.Lock()

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
// args[0] 为缓存值 args[1]如果存在, 则覆盖默认过期时间, 单位秒
func (cm *CacheManager) SetNX(clib string, key string, args ...any) (bool, error) {
	defer cm.locker.Unlock()
	cm.locker.Lock()

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
// args[0] 为缓存值 args[1]如果存在, 则覆盖默认过期时间, 单位秒
func (cm *CacheManager) Incrby(clib string, key string, args ...any) (int64, error) {
	defer cm.locker.Unlock()
	cm.locker.Lock()

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

// Del 删除缓存信息
func (cm *CacheManager) Del(clib string, key string) error {
	defer cm.locker.Unlock()
	cm.locker.Lock()

	if tm, ok := cm.clibs[clib]; !ok {
		return ipakku.ErrCacheLibNotExist
	} else {
		tm.DestroyToken(key)
	}
	return nil
}

// Clear 清空库内容
func (cm *CacheManager) Clear(clib string) {
	defer cm.locker.Unlock()
	cm.locker.Lock()

	if tm, ok := cm.clibs[clib]; ok {
		tm.Clear()
	}
}

// getExpSecond 获取过期时间, 若args[1]有值, 则返回args[1]的值, 否则返回之前注册lib时的值
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
