// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package ipakku

import "errors"

// ErrCacheLibNotExist 缓存库没有注册
var ErrCacheLibNotExist = errors.New("cache lib not exist")

// ErrCacheLibIsExist 缓存库重复注册
var ErrCacheLibIsExist = errors.New("cache lib is exist")

// ErrNoCacheHit 没有命中缓存
var ErrNoCacheHit = errors.New("no cache hit")

// AppCache app 缓存模块
type AppCache interface {

	// lib为库名, second:过期时间-1为不过期
	RegLib(clib string, second int64) error

	// Set 向lib库中设置键为key的值tb
	Set(clib string, key string, tb interface{}) error

	// Get 读取缓存信息
	Get(clib string, key string, val interface{}) error

	// Del 删除缓存信息
	Del(clib string, key string) error

	// Keys 获取库的所有key
	Keys(clib string) []string

	// Clear 清空库内容
	Clear(clib string)
}
