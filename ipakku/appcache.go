// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package ipakku

import "errors"

// ErrCacheLibNotExist 缓存库没有注册
var ErrCacheLibNotExist = errors.New("cache lib not exist")

// ErrCacheLibIsExist 缓存库重复注册
var ErrCacheLibIsExist = errors.New("cache lib is exist")

// ErrNoCacheHit 没有命中缓存
var ErrNoCacheHit = errors.New("no cache hit")

// ErrCacheArgsEmpty 必填字段为空
var ErrCacheArgsEmpty = errors.New("cache value parameter cannot be empty")

// ErrCacheArgsTypeError 缓存值参数类型错误
var ErrCacheArgsTypeError = errors.New("cache parameter type error")

// ErrCacheConvertError 缓存参数类型错误
var ErrCacheConvertError = errors.New("cache parameter type error")

// AppCache 缓存模块
type AppCache interface {

	// RegLib lib: 库名(组名), second: 默认过期时间, -1为不过期
	RegLib(clib string, second int64) error

	// Exists 返回key是否存在
	Exists(clib string, key string) (bool, error)

	// Set 向lib库中设置键为key的值
	// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
	Set(clib string, key string, args ...any) error

	// SetNX 向lib库中设置键为key的值, 当key不存在时设置成功, 并返回true
	// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
	SetNX(clib string, key string, args ...any) (bool, error)

	// Incrby 指定key以increment的值累加, 返回累加后的值
	// args[0] 为缓存值 args[2]如果存在, 则覆盖默认过期时间, 单位秒
	Incrby(clib string, key string, args ...any) (int64, error)

	// Get 读取缓存信息
	Get(clib string, key string, val any) error

	// Del 删除缓存信息
	Del(clib string, key string) error

	// Keys 获取库的所有key
	Keys(clib string) []string

	// Clear 清空库内容
	Clear(clib string)
}

// ICache 缓存接口
type ICache interface {
	AppCache

	// Init 初始化缓存管理器, 一个对象只能初始化一次
	Init(config AppConfig, appName string)
}
