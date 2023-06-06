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

const (
	flagIndir uintptr = 1 << 7
)

type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
	flag uintptr
}
type rtype struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32
	tflag      uint8
	align      uint8
	fieldAlign uint8
	kind       uint8
	equal      func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata     *byte
	str        int32
	ptrToThis  int32
}

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(t *rtype, dst, src unsafe.Pointer)

//go:linkname typedmemclr reflect.typedmemclr
func typedmemclr(t *rtype, ptr unsafe.Pointer)

//go:linkname assignTo reflect.(*Value).assignTo
func assignTo(v *reflect.Value, context string, dst *rtype, target unsafe.Pointer) reflect.Value

//go:linkname zeroVal reflect.zeroVal
var zeroVal [1024]byte

// GetFunctionName 获取函数名称
func GetFunctionName(i interface{}, seps ...rune) string {
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
func GetTagValues(tagName string, obj interface{}) (rs map[string]string, err error) {
	rs = make(map[string]string)
	var fields []reflect.StructField
	if fields, err = GetTagField(obj); nil == err && len(fields) > 0 {
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
	return rs, err
}

// GetTagFieldName 获取结构体, 含有tagName的字段
func GetTagFieldName(tagName string, ptr interface{}) (rs []string, err error) {
	var fields []reflect.StructField
	if fields, err = GetTagField(ptr); nil == err && len(fields) > 0 {
		for i := 0; i < len(fields); i++ {
			if s := string(fields[i].Tag); s == tagName || strings.HasPrefix(s, tagName) {
				rs = append(rs, fields[i].Name)
			}
		}
	}
	return
}

// GetTagField 获取结构体的字段
func GetTagField(obj interface{}) (res []reflect.StructField, err error) {
	t := GetNotPtrRefType(obj)
	for i := 0; i < t.NumField(); i++ {
		res = append(res, t.Field(i))
	}
	return res, nil
}

// GetNotPtrRefType 获取结构体的字段类型
func GetNotPtrRefType(obj interface{}) reflect.Type {
	var t reflect.Type
	if v, ok := obj.(reflect.Type); ok {
		t = v
	} else if v, ok := obj.(reflect.Value); ok {
		t = v.Type()
	} else {
		t = reflect.TypeOf(obj)
	}
	if t.Kind() == reflect.Ptr {
		for t = t.Elem(); t != nil && t.Kind() == reflect.Ptr; {
			if t = t.Elem(); nil == t {
				return t
			}
		}
	}
	return t
}

// GetStructFieldType 获取结构体的类型
func GetStructFieldType(obj interface{}, fieldName string) (reflect.Type, error) {
	t := GetNotPtrRefType(obj)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("only struct are supported")
	}
	if f, ok := t.FieldByName(fieldName); ok {
		return f.Type, nil
	}
	return nil, nil
}

// GetStructFieldRefValue 获取结构体的值
func GetStructFieldRefValue(src interface{}, fieldName string) (reflect.Value, error) {
	t := reflect.TypeOf(src)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("only pointer struct are supported")
	}
	return reflect.ValueOf(src).Elem().FieldByName(fieldName), nil
}

// SetStructFieldValue 将结构体里的成员按照json名字来赋值
func SetStructFieldValue(dstStruct interface{}, fieldName string, val interface{}) error {
	t := reflect.TypeOf(dstStruct)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer struct are supported")
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

// SetStructFieldValueUnSafe 给结构体里的成员字段赋值 - 非安全指针, 可以设置私有值
func SetStructFieldValueUnSafe(dstStruct interface{}, targetField string, obj interface{}) error {
	if st := reflect.TypeOf(dstStruct); st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer struct are supported")
	}
	valueOfTargetField := reflect.ValueOf(dstStruct).Elem().FieldByName(targetField)
	reflect.NewAt(valueOfTargetField.Type(), unsafe.Pointer(valueOfTargetField.UnsafeAddr())).Elem().Set(reflect.ValueOf(obj))
	return nil
}

// SetInterfaceValueUnSafe 给接口类型的src赋值val - 非安全指针
func SetInterfaceValueUnSafe(dst interface{}, val interface{}) error {
	if st := reflect.TypeOf(dst); st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Interface {
		return errors.New("only pointer interface are supported")
	}
	valueOfTarget := reflect.ValueOf(dst).Elem()
	reflect.NewAt(valueOfTarget.Type(), unsafe.Pointer(valueOfTarget.UnsafeAddr())).Elem().Set(reflect.ValueOf(val))
	return nil
}

// SetStructFieldValueUnSafe 给结构体里的成员字段赋值 - 非安全指针, 可以设置私有值
// 参考reflect.Set函数改造
// func SetStructFieldValueUnSafe1(dstStruct interface{}, targetField string, obj interface{}) error {
// 	// 仅支持指针类型结构体
// 	if st := reflect.TypeOf(dstStruct); st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Struct {
// 		return errors.New("only pointer struct are supported")
// 	}
// 	// 转换obj对象为value对象
// 	var valueOfObj *reflect.Value
// 	if convertedVal, ok := obj.(*reflect.Value); ok {
// 		elem := convertedVal.Elem()
// 		valueOfObj = &elem
// 	} else {
// 		if vt := reflect.TypeOf(obj); vt.Kind() != reflect.Ptr {
// 			return errors.New("only pointer object are supported")
// 		}
// 		v := reflect.ValueOf(obj)
// 		valueOfObj = &v
// 	}
// 	// 获取目标字段的指针
// 	valueOfTargetField := reflect.ValueOf(dstStruct).Elem().FieldByName(targetField)
// 	targetPointer := (*emptyInterface)(unsafe.Pointer(&valueOfTargetField))
// 	var target unsafe.Pointer
// 	if valueOfTargetField.Kind() == reflect.Interface {
// 		target = targetPointer.word
// 	}
// 	// 加权限
// 	valueOfAssign := assignTo(valueOfObj, "reflect.Set", targetPointer.typ, target)
// 	valuePointer := (*emptyInterface)(unsafe.Pointer(&valueOfAssign))
// 	// 改变数据
// 	if valuePointer.flag&flagIndir != 0 {
// 		if valuePointer.word == unsafe.Pointer(&zeroVal[0]) {
// 			typedmemclr(targetPointer.typ, targetPointer.word)
// 		} else {
// 			typedmemmove(targetPointer.typ, targetPointer.word, valuePointer.word)
// 		}
// 	} else {
// 		*(*unsafe.Pointer)(targetPointer.word) = valuePointer.word
// 	}
// 	return nil
// }

// // SetInterfaceValueUnSafe 给接口类型的src赋值val - 非安全指针
// func SetInterfaceValueUnSafe1(src interface{}, val interface{}) error {
// 	if st := reflect.TypeOf(src); st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Interface {
// 		return errors.New("only pointer interface are supported")
// 	}
// 	var target unsafe.Pointer
// 	sv := reflect.ValueOf(src).Elem()
// 	spointer := (*emptyInterface)(unsafe.Pointer(&sv))
// 	if tv := reflect.TypeOf(val); tv.Kind() != reflect.Ptr {
// 		return errors.New("only pointer objects are supported")
// 	} else {
// 		if tv.Elem().Kind() == reflect.Interface {
// 			target = spointer.word
// 		}
// 	}
// 	vv := reflect.ValueOf(val)
// 	vvx := assignTo(&vv, "reflect.Set", spointer.typ, target)
// 	vpointer := (*emptyInterface)(unsafe.Pointer(&vvx))
// 	//
// 	if vpointer.flag&flagIndir != 0 {
// 		typedmemmove(spointer.typ, spointer.word, vpointer.word)
// 	} else {
// 		*(*unsafe.Pointer)(spointer.word) = vpointer.word
// 	}
// 	return nil
// }

// Invoke 调用src里的方法, 返回 []reflect.Value
func Invoke(src interface{}, method string, params ...interface{}) []reflect.Value {
	args := make([]reflect.Value, len(params))
	if len(params) > 0 {
		for i, temp := range params {
			args[i] = reflect.ValueOf(temp)
		}
	}
	return reflect.ValueOf(src).MethodByName(method).Call(args)
}
