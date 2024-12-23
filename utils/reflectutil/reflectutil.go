// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 反射工具

package reflectutil

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

// GetFunctionName 获取函数名称
func GetFunctionName(i any, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

// GetTagValues 获取结构体, 含有tagName的字段和值
func GetTagValues(tagName string, obj any) (rs map[string]string) {
	rs = make(map[string]string)
	if fields := GetStructFields(obj); len(fields) > 0 {
		for i := 0; i < len(fields); i++ {
			if s := string(fields[i].Tag); s == tagName {
				rs[fields[i].Name] = ""
			} else if strings.HasPrefix(s, tagName) {
				if tagVal := fields[i].Tag.Get(tagName); len(tagVal) > 0 {
					rs[fields[i].Name] = tagVal
				} else {
					rs[fields[i].Name] = ""
				}
			}
		}
	}
	return rs
}

// GetTagFieldName 获取结构体, 含有tagName的字段
func GetTagFieldName(tagName string, ptr any) (rs []string) {
	if fields := GetStructFields(ptr); len(fields) > 0 {
		for i := 0; i < len(fields); i++ {
			if s := string(fields[i].Tag); s == tagName || strings.HasPrefix(s, tagName) {
				rs = append(rs, fields[i].Name)
			}
		}
	}
	return
}

// GetAnonymousField 获取匿名结构体字段
func GetAnonymousField(obj any) (res []reflect.StructField) {
	if fields := GetStructFields(obj); len(fields) > 0 {
		for i := 0; i < len(fields); i++ {
			if fields[i].Anonymous {
				res = append(res, fields[i])
			}
		}
	}
	return res
}

// GetAnonymousOrNoneTypeNameField 获取匿名和未命名嵌套结构体字段
func GetAnonymousOrNoneTypeNameField(obj any) (res []reflect.StructField) {
	if fields := GetStructFields(obj); len(fields) > 0 {
		for i := 0; i < len(fields); i++ {
			if !fields[i].Anonymous {
				if tp := GetNotPtrRefType(fields[i].Type); reflect.Struct != tp.Kind() || len(tp.Name()) > 0 {
					continue
				}
			}
			res = append(res, fields[i])
		}
	}
	return res
}

// GetStructFields 获取结构体的字段
func GetStructFields(obj any) (res []reflect.StructField) {
	t := GetNotPtrRefType(obj)
	for i := 0; i < t.NumField(); i++ {
		res = append(res, t.Field(i))
	}
	return res
}

// GetNotPtrRefType 获取结构体的字段类型
func GetNotPtrRefType(obj any) reflect.Type {
	var t reflect.Type
	if v, ok := obj.(reflect.Type); ok {
		t = v
	} else if v, ok := obj.(reflect.Value); ok {
		t = v.Type()
	} else {
		t = reflect.TypeOf(obj)
	}
	if nil != t && t.Kind() == reflect.Ptr {
		for t = t.Elem(); t != nil && t.Kind() == reflect.Ptr; {
			if t = t.Elem(); nil == t {
				return t
			}
		}
	}
	return t
}

// GetStructFieldType 获取结构体的类型
func GetStructFieldType(obj any, fieldName string) (reflect.Type, error) {
	t := GetNotPtrRefType(obj)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("only struct are supported, but input type is: " + t.Kind().String())
	}
	if f, ok := t.FieldByName(fieldName); ok {
		return f.Type, nil
	}
	return nil, nil
}

// GetStructFieldRefValue 获取结构体的值
func GetStructFieldRefValue(src any, fieldName string) (reflect.Value, error) {
	if err := assertionObjectType(src, true, reflect.Struct); nil != err {
		return reflect.Value{}, err
	}
	if srcval, ok := src.(reflect.Value); ok {
		return srcval.Elem().FieldByName(fieldName), nil
	}
	return reflect.ValueOf(src).Elem().FieldByName(fieldName), nil
}

// SetStructFieldValue 给结构体里内指定的成员变量赋值
func SetStructFieldValue(dstStruct any, fieldName string, val any) error {
	if err := assertionObjectType(dstStruct, true, reflect.Struct); nil != err {
		return err
	}
	v := reflect.ValueOf(dstStruct).Elem()
	field := v.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		if reflect.ValueOf(val).Type().AssignableTo(field.Type()) {
			field.Set(reflect.ValueOf(val))
			return nil
		}
	}
	return errors.New("value of type " + reflect.ValueOf(val).Type().String() + " is not assignable to type " + field.Type().String())
}

// SetStructFieldValueUnSafe 给结构体里的成员字段赋值 - 可以设置私有值
func SetStructFieldValueUnSafe(dstStruct any, targetField string, obj any) error {
	if err := assertionObjectType(dstStruct, true, reflect.Struct); nil != err {
		return err
	}
	var valueOfTargetField reflect.Value
	if dstStructVal, ok := dstStruct.(reflect.Value); ok {
		valueOfTargetField = dstStructVal.Elem().FieldByName(targetField)
	} else {
		valueOfTargetField = reflect.ValueOf(dstStruct).Elem().FieldByName(targetField)
	}
	reflect.NewAt(valueOfTargetField.Type(), unsafe.Pointer(valueOfTargetField.UnsafeAddr())).Elem().Set(reflect.ValueOf(obj))
	return nil
}

// SetInterfaceValueUnSafe 给接口类型的src赋值val
func SetInterfaceValueUnSafe(dst any, val any) error {
	if err := assertionObjectType(dst, true, reflect.Interface); nil != err {
		return err
	}
	var valueOfTarget reflect.Value
	if dstStructVal, ok := dst.(reflect.Value); ok {
		valueOfTarget = dstStructVal.Elem()
	} else {
		valueOfTarget = reflect.ValueOf(dst).Elem()
	}
	reflect.NewAt(valueOfTarget.Type(), unsafe.Pointer(valueOfTarget.UnsafeAddr())).Elem().Set(reflect.ValueOf(val))
	return nil
}

// Invoke 调用src里的方法, 返回 []reflect.Value
func Invoke(src any, method string, params ...any) []reflect.Value {
	args := make([]reflect.Value, len(params))
	if len(params) > 0 {
		for i, temp := range params {
			args[i] = reflect.ValueOf(temp)
		}
	}
	return reflect.ValueOf(src).MethodByName(method).Call(args)
}

// assertionObjectType 判断输入对象类型
func assertionObjectType(inputObj any, isPointer bool, types ...reflect.Kind) error {
	var st reflect.Type
	if reftp, ok := inputObj.(reflect.Type); ok {
		st = reftp
	} else if refval, ok := inputObj.(reflect.Value); ok {
		st = refval.Type()
	} else if st = reflect.TypeOf(inputObj); nil == st {
		return errors.New("input object is nil")
	}

	if isPointer && st.Kind() != reflect.Ptr {
		return errors.New("only pointer object are supported, but the input type is: " + st.Kind().String())
	}

	if lenTypes := len(types); lenTypes > 0 {
		var stKind reflect.Kind
		if isPointer {
			stKind = st.Elem().Kind()
		} else {
			stKind = st.Kind()
		}

		has := false
		if lenTypes == 1 {
			has = stKind == types[0]
		} else {
			for i := 0; i < lenTypes; i++ {
				if stKind == types[i] {
					has = true
					break
				}
			}
		}
		if !has {
			return errors.New("the current input objec type [" + stKind.String() + "] is not supported")
		}
	}

	return nil
}
