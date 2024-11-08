// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 拓展对象-带安全锁的map

package utypes

import (
	"errors"
	"sync"
)

// NewSafeMap 新建带RWMutex锁的map
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		lock: new(sync.RWMutex),
		cmap: make(map[K]V),
	}
}

// SafeMap 带RWMutex锁的map
type SafeMap[K comparable, V any] struct {
	lock *sync.RWMutex
	cmap map[K]V
}

// New 初始化
func (m SafeMap[K, V]) New() *SafeMap[K, V] {
	r := &m
	r.lock = new(sync.RWMutex)
	r.cmap = make(map[K]V)
	return r
}

// Get 获取值
func (m *SafeMap[K, V]) Get(k K) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	val, ok := m.cmap[k]
	return val, ok
}

// Cut 获取值, 剪切方式
func (m *SafeMap[K, V]) Cut(k K) (V, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	val, ok := m.cmap[k]
	if ok {
		delete(m.cmap, k)
	}
	return val, ok
}

// CutR 随机获取一个值, 剪切方式
func (m *SafeMap[K, V]) CutR() (res V, exist bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if len(m.cmap) == 0 {
		return
	}
	var key K
	for key = range m.cmap {
		break
	}
	if res, exist = m.cmap[key]; exist {
		delete(m.cmap, key)
	}
	return
}

// Put 插入值
func (m *SafeMap[K, V]) Put(k K, v V) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cmap[k] = v
}

// Put 插入值, 如果存在则报错
func (m *SafeMap[K, V]) PutX(k K, v V) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.cmap[k]; ok {
		return errors.New("key is exist")
	}
	m.cmap[k] = v
	return nil
}

// Keys 获取所有的key
func (m *SafeMap[K, V]) Keys() []K {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make([]K, len(m.cmap))
	i := 0
	for k := range m.cmap {
		r[i] = k
		i++
	}
	return r
}

// Values 获取所有的value
func (m *SafeMap[K, V]) Values() []V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make([]V, len(m.cmap))
	i := 0
	for _, val := range m.cmap {
		r[i] = val
		i++
	}
	return r
}

// ContainsKey  是否包含key
func (m *SafeMap[K, V]) ContainsKey(k K) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.cmap[k]
	return ok
}

// Delete 删除
func (m *SafeMap[K, V]) Delete(k K) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.cmap, k)
}

// Clear 清空
func (m *SafeMap[K, V]) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.cmap = make(map[K]V)
}

// ToMap 获取map值, 复制值
func (m *SafeMap[K, V]) ToMap() map[K]V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[K]V)
	for k, v := range m.cmap {
		r[k] = v
	}
	return r
}

// DoRange 循环
func (m *SafeMap[K, V]) DoRange(fun func(key K, val V) error) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for k, v := range m.cmap {
		if err := fun(k, v); nil != err {
			return err
		}
	}
	return nil
}

// Size 返回大小
func (m *SafeMap[K, V]) Size() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.cmap)
}
