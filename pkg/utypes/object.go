// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 拓展对象-any转各种类型

package utypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Object interface类型转换
type Object struct {
	// o 实例化时保存的原对象或指针
	o any
}

// NewObject 新建一个object对象
func NewObject(obj any) Object {
	return Object{o: obj}
}

// ToBool 转换为bool
func (obj Object) ToBool(d bool) bool {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(bool); ok {
		return r
	} else if r, ok := obj.o.(string); ok {
		return strings.ToUpper(r) == "TRUE"
	}
	return d
}

// ToString 转换为string
func (obj Object) ToString(d string) string {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(string); ok && len(r) > 0 {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		return r.String()
	}
	return d
}

// ToByte 转换为string
func (obj Object) ToByte(d []byte) []byte {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.([]byte); ok {
		return r
	}
	return d
}

// ToInt 转换为int
func (obj Object) ToInt(d int) int {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(int); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return int(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseInt(r, 10, 64); nil == err {
			return int(r)
		}
		return d
	}
	return d
}

// ToInt8 转换为int8
func (obj Object) ToInt8(d int8) int8 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(int8); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return int8(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseInt(r, 10, 64); nil == err {
			return int8(r)
		}
		return d
	}
	return d
}

// ToInt16 转换为int16
func (obj Object) ToInt16(d int16) int16 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(int16); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return int16(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseInt(r, 10, 64); nil == err {
			return int16(r)
		}
		return d
	}
	return d
}

// ToInt32 转换为int32
func (obj Object) ToInt32(d int32) int32 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(int32); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return int32(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseInt(r, 10, 64); nil == err {
			return int32(r)
		}
		return d
	}
	return d
}

// ToInt64 转换为int64
func (obj Object) ToInt64(d int64) int64 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(int64); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return v
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseInt(r, 10, 64); nil == err {
			return r
		}
		return d
	}
	return d
}

// ToUint 转换为uint
func (obj Object) ToUint(d uint) uint {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(uint); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return uint(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseUint(r, 10, 64); nil == err {
			return uint(r)
		}
		return d
	}
	return d
}

// ToUint8 转换为uint8
func (obj Object) ToUint8(d uint8) uint8 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(uint8); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return uint8(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseUint(r, 10, 64); nil == err {
			return uint8(r)
		}
		return d
	}
	return d
}

// ToUint16 转换为uint16
func (obj Object) ToUint16(d uint16) uint16 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(uint16); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return uint16(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseUint(r, 10, 64); nil == err {
			return uint16(r)
		}
		return d
	}
	return d
}

// ToUint32 转换为uint32
func (obj Object) ToUint32(d uint32) uint32 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(uint32); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return uint32(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseUint(r, 10, 64); nil == err {
			return uint32(r)
		}
		return d
	}
	return d
}

// ToUint64 转换为uint64
func (obj Object) ToUint64(d uint64) uint64 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(uint64); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Int64(); nil == err {
			return uint64(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseUint(r, 10, 64); nil == err {
			return r
		}
		return d
	}
	return d
}

// ToFloat32 转换为float32
func (obj Object) ToFloat32(d float32) float32 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(float32); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Float64(); nil == err {
			return float32(v)
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseFloat(r, 64); nil == err {
			return float32(r)
		}
		return d
	}
	return d
}

// ToFloat64 转换为Float64
func (obj Object) ToFloat64(d float64) float64 {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(float64); ok {
		return r
	} else if r, ok := obj.o.(json.Number); ok {
		if v, err := r.Float64(); nil == err {
			return v
		}
	} else if r, ok := obj.o.(string); ok {
		if r, err := strconv.ParseFloat(r, 64); nil == err {
			return r
		}
		return d
	}
	return d
}

// ToStrMap 转换为map[string]any
func (obj Object) ToStrMap(d map[string]any) map[string]any {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(map[string]any); ok {
		return r
	}
	return d
}

// ToIntMap 转换为map[int]any
func (obj Object) ToIntMap(d map[int]any) map[int]any {
	if obj.IsNill() {
		return d
	}

	if r, ok := obj.o.(map[int]any); ok {
		return r
	}
	return d
}

// ToInt32Map 转换为map[int32]any
func (obj Object) ToInt32Map(d map[int32]any) map[int32]any {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(map[int32]any); ok {
		return r
	}
	return d
}

// ToInt64Map 转换为map[int64]any
func (obj Object) ToInt64Map(d map[int64]any) map[int64]any {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(map[int64]any); ok {
		return r
	}
	return d
}

// ToFloat32Map 转换为map[float32]any
func (obj Object) ToFloat32Map(d map[float32]any) map[float32]any {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(map[float32]any); ok {
		return r
	}
	return d
}

// ToFloat64Map 转换为map[float64]any
func (obj Object) ToFloat64Map(d map[float64]any) map[float64]any {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(map[float64]any); ok {
		return r
	}
	return d
}

// IsNill 是否是空
func (obj Object) IsNill() bool {
	return nil == obj.o
}

// SetVal 设置原始值
func (obj *Object) SetVal(newVal any) *Object {
	obj.o = newVal
	return obj
}

// GetVal 获取原始值
func (obj Object) GetVal() any {
	return obj.o
}

// Scan 自动赋值, 如果类型不一致会尝试转换
func (obj Object) Scan(v any) error {
	if obj.IsNill() {
		return nil // 空值不处理
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("non-pointer or nil passed to Scan")
	}

	// 获取指针指向的实际值
	target := rv.Elem()
	srcValue := reflect.ValueOf(obj.o)

	// 类型完全匹配时直接赋值
	if srcValue.Type().AssignableTo(target.Type()) {
		target.Set(srcValue)
		return nil
	}

	// 处理基础类型转换
	switch target.Kind() {
	case reflect.Bool:
		target.SetBool(obj.ToBool(false))
		return nil

	case reflect.String:
		target.SetString(obj.ToString(""))
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal := obj.ToInt64(0)
		if target.OverflowInt(intVal) {
			return fmt.Errorf("value %d overflows %s", intVal, target.Type())
		}
		target.SetInt(intVal)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal := obj.ToUint64(0)
		if target.OverflowUint(uintVal) {
			return fmt.Errorf("value %d overflows %s", uintVal, target.Type())
		}
		target.SetUint(uintVal)
		return nil

	case reflect.Float32, reflect.Float64:
		floatVal := obj.ToFloat64(0)
		if target.Kind() == reflect.Float32 && (floatVal > math.MaxFloat32 || floatVal < -math.MaxFloat32) {
			return fmt.Errorf("value %f overflows float32", floatVal)
		}
		target.SetFloat(floatVal)
		return nil
	}

	// 最后尝试JSON序列化/反序列化
	jsonBytes, err := json.Marshal(obj.o)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, v); err != nil {
		return fmt.Errorf("json unmarshal error: %w", err)
	}

	return nil
}
