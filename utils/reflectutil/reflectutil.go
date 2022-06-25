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

//go:linkname assignTo reflect.(*Value).assignTo
func assignTo(v *reflect.Value, context string, dst *rtype, target unsafe.Pointer) reflect.Value

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
func GetTagValues(tagName string, ptr interface{}) (map[string]string, error) {
	var v reflect.Value
	t := reflect.TypeOf(ptr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		v = reflect.ValueOf(ptr)
	} else {
		v = reflect.ValueOf(ptr).Elem()
	}
	rs := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tagVal := field.Tag.Get(tagName)
		if len(tagVal) > 0 {
			rs[field.Name] = tagVal
		}
	}
	return rs, nil
}

// GetStructFieldType 获取结构体的类型
func GetStructFieldType(ptr interface{}, fieldName string) (reflect.Type, error) {
	t := reflect.TypeOf(ptr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, errors.New("only pointer objects are supported")
	}
	v := reflect.ValueOf(ptr).Elem()
	field := v.FieldByName(fieldName)
	return field.Type(), nil
}

// SetStructFieldValue 将结构体里的成员按照json名字来赋值
func SetStructFieldValue(ptr interface{}, fieldName string, val interface{}) error {
	t := reflect.TypeOf(ptr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer objects are supported")
	}
	v := reflect.ValueOf(ptr).Elem()
	field := v.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		if reflect.ValueOf(val).Type().AssignableTo(field.Type()) {
			field.Set(reflect.ValueOf(val))
			return nil
		}
	}
	return errors.New("value of type " + reflect.ValueOf(val).Type().String() + " is not assignable to type " + field.Type().String())
}

// SetStructFieldValueUnSafe 将结构体里的成员按照json名字来赋值 - 非安全指针, 可以设置私有值
func SetStructFieldValueUnSafe(src interface{}, fieldName string, val interface{}) error {
	if st := reflect.TypeOf(src); st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer struct are supported")
	}
	if vt := reflect.TypeOf(val); vt.Kind() != reflect.Ptr {
		return errors.New("only pointer object are supported")
	}
	//
	fieldObj := reflect.ValueOf(src).Elem().FieldByName(fieldName)
	spointer := (*emptyInterface)(unsafe.Pointer(&fieldObj))
	var target unsafe.Pointer
	if fieldObj.Kind() == reflect.Interface {
		target = spointer.word
	}
	vv := reflect.ValueOf(val)
	vvx := assignTo(&vv, "reflect.Set", spointer.typ, target)
	vpointer := (*emptyInterface)(unsafe.Pointer(&vvx))
	//
	if vpointer.flag&flagIndir != 0 {
		typedmemmove(spointer.typ, spointer.word, vpointer.word)
	} else {
		*(*unsafe.Pointer)(spointer.word) = vpointer.word
	}
	return nil
}

// SetInterfaceValueUnSafe 给接口类型的src赋值val - 非安全指针
func SetInterfaceValueUnSafe(src interface{}, val interface{}) error {
	if st := reflect.TypeOf(src); st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Interface {
		return errors.New("only pointer interface are supported")
	}
	var target unsafe.Pointer
	sv := reflect.ValueOf(src).Elem()
	spointer := (*emptyInterface)(unsafe.Pointer(&sv))
	if tv := reflect.TypeOf(val); tv.Kind() != reflect.Ptr {
		return errors.New("only pointer objects are supported")
	} else {
		if tv.Elem().Kind() == reflect.Interface {
			target = spointer.word
		}
	}
	vv := reflect.ValueOf(val)
	vvx := assignTo(&vv, "reflect.Set", spointer.typ, target)
	vpointer := (*emptyInterface)(unsafe.Pointer(&vvx))
	//
	if vpointer.flag&flagIndir != 0 {
		typedmemmove(spointer.typ, spointer.word, vpointer.word)
	} else {
		*(*unsafe.Pointer)(spointer.word) = vpointer.word
	}
	return nil
}

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
