// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 拓展对象-any转各种类型

package utypes

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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

// ToTime 转换为time
func (obj Object) ToTime(d time.Time) time.Time {
	if obj.IsNill() {
		return d
	}
	if r, ok := obj.o.(time.Time); ok {
		return r
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

// Scan 自动赋值
func (obj Object) Scan(v any) error {
	switch v := v.(type) {
	case nil:
		return nil
	case *string:
		if !obj.IsNill() {
			*v = obj.ToString("")
		}
		return nil
	case *[]byte:
		if !obj.IsNill() {
			*v = obj.ToByte(nil)
		}
		return nil
	case *int:
		if !obj.IsNill() {
			*v = obj.ToInt(0)
		}
		return nil
	case *int8:
		if !obj.IsNill() {
			*v = obj.ToInt8(0)
		}
		return nil
	case *int16:
		if !obj.IsNill() {
			*v = obj.ToInt16(0)
		}
		return nil
	case *int32:
		if !obj.IsNill() {
			*v = obj.ToInt32(0)
		}
		return nil
	case *int64:
		if !obj.IsNill() {
			*v = obj.ToInt64(0)
		}
		return nil
	case *uint:
		if !obj.IsNill() {
			*v = obj.ToUint(0)
		}
		return nil
	case *uint8:
		if !obj.IsNill() {
			*v = obj.ToUint8(0)
		}
		return nil
	case *uint16:
		if !obj.IsNill() {
			*v = obj.ToUint16(0)
		}
		return nil
	case *uint32:
		if !obj.IsNill() {
			*v = obj.ToUint32(0)
		}
		return nil
	case *uint64:
		if !obj.IsNill() {
			*v = obj.ToUint64(0)
		}
		return nil
	case *float32:
		if !obj.IsNill() {
			*v = obj.ToFloat32(0)
		}
		return nil
	case *float64:
		if !obj.IsNill() {
			*v = obj.ToFloat64(0)
		}
		return nil
	case *bool:
		if !obj.IsNill() {
			*v = obj.ToBool(false)
		}
		return nil
	case *time.Time:
		if !obj.IsNill() {
			*v = obj.ToTime(time.Time{})
		}
		return nil
	case encoding.BinaryUnmarshaler:
		if !obj.IsNill() {
			if ms, ok := obj.o.(encoding.BinaryMarshaler); ok {
				if val, err := ms.MarshalBinary(); nil == err {
					return v.UnmarshalBinary(val)
				} else {
					return err
				}
			} else {
				return errors.New("object not support encoding.BinaryMarshaler")
			}
		}
		return nil
	default:
		return fmt.Errorf("can't support scan %T (consider implementing BinaryMarshaler and BinaryUnmarshaler)", v)
	}
}
