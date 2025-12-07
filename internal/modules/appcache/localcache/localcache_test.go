// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package localcache

import (
	"strconv"
	"testing"
	"time"
)

// 缓存模块
func TestCacheManager(t *testing.T) {
	clib := "_test_"
	exp := int64(5)
	count := 100000
	cachemanager := &CacheManager{}
	cachemanager.Init(nil, "")
	err := cachemanager.RegLib(clib, exp)
	if nil != err {
		t.Log(err)
		t.FailNow()
	}
	// 赋值
	for i := 0; i < count; i++ {
		key := "key_" + strconv.FormatInt(int64(i), 10)
		err = cachemanager.Set(clib, key, i)
		if nil != err {
			t.Log(err)
			t.FailNow()
		}
		var val int
		if err := cachemanager.Get(clib, key, &val); err != nil {
			t.Error(err)
		} else {
			t.Log(val)
		}
	}
	keys := cachemanager.Keys(clib)
	if len(keys) == 0 {
		t.Log("无法列出所有key")
		t.FailNow()
	}
	// 取值
	time.Sleep(time.Duration(5) * time.Second)
	for i := 0; i < count; i++ {
		key := "key_" + strconv.FormatInt(int64(i), 10)
		var val int
		if err := cachemanager.Get(clib, key, &val); nil == err {
			t.Log("缓存未销毁")
			break
		}
	}
}

// 缓存模块
func BenchmarkCacheManager(t *testing.B) {
	clib := "_test_"
	exp := int64(5)
	count := 5
	var err error
	cachemanager := &CacheManager{}
	cachemanager.Init(nil, "")
	for i := 0; i < count; i++ {
		si := strconv.FormatInt(int64(i), 10)
		err = cachemanager.RegLib(clib+si, exp+int64(i))
		if nil != err {
			t.Log(err)
			t.FailNow()
		}
	}
	// 赋值
	for i := 0; i < count; i++ {
		si := strconv.FormatInt(int64(i), 10)
		err = cachemanager.Set(clib+si, si, i)
		if nil != err {
			t.Log(err)
			t.FailNow()
		}
		var val int
		if err := cachemanager.Get(clib+si, si, &val); nil != err {
			t.Log("无法查询到缓存")
			t.Error(err)
			break
		}
		keys := cachemanager.Keys(clib + si)
		if len(keys) == 0 {
			t.Log("无法列出所有key")
			t.FailNow()
		}
	}
}
