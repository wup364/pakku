// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package localcache

import (
	"github.com/wup364/pakku/ipakku"
)

// StructValue 结构体的值, 包含一个任意类型的Value
type StructValue struct {
	Value any
}

// LocalCacheValueScan 缓存值对象转换接口, 缓存值若实现此接口, Get时会调用
func (sv *StructValue) LocalCacheValueScan(val any) error {
	if uat, ok := val.(*StructValue); ok {
		uat.Value = sv.Value
		return nil
	}
	return ipakku.ErrCacheConvertError
}
