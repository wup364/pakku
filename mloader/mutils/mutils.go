// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package mutils

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/reflectutil"
)

// AutoWired 自动注入依赖
func AutoWired(ptr interface{}, l ipakku.Loader) (err error) {
	var tagvals = make(map[string]string)
	if tagvals = reflectutil.GetTagValues(ipakku.PAKKUTAG_AUTOWIRED, ptr); len(tagvals) == 0 {
		return autoWiredAnonymousStruct(ptr, l)
	}

	// 仅支持指针类型结构体
	if t := reflect.TypeOf(ptr); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer objects are supported")
	}
	for field, valKey := range tagvals {
		var val interface{}
		if val, err = getModuleByName(valKey, l); nil != err {
			break
		}
		var ftype reflect.Type
		if ftype, err = reflectutil.GetStructFieldType(ptr, field); nil == err {
			logs.Infof("> Autowired %s <= %s[%s] \r\n", field, valKey, ftype.String())
			if ftype.Kind() != reflect.Interface {
				err = fmt.Errorf("only interface type injections are accepted, field: %s", field)
				break
			}
			if err = reflectutil.SetStructFieldValueUnSafe(ptr, field, val); nil != err {
				break
			}
		} else {
			logs.Infof("> Autowired %s <= %s[err=%s] \r\n", field, valKey, err.Error())
		}
	}

	if nil == err {
		// 自动注入匿名嵌套结构体
		err = autoWiredAnonymousStruct(ptr, l)
	}
	return err
}

// autoWiredAnonymousStruct 自动注入匿名嵌套结构体
func autoWiredAnonymousStruct(ptr interface{}, l ipakku.Loader) (err error) {
	var fields []reflect.StructField
	if fields = reflectutil.GetAnonymousField(ptr); len(fields) == 0 {
		return
	}

	for i := 0; i < len(fields); i++ {
		// 仅支持指针类型结构体
		if t := fields[i].Type; t.Kind() != reflect.Struct {
			continue
		}
		var refval reflect.Value
		if refval, err = reflectutil.GetStructFieldRefValue(ptr, fields[i].Name); nil != err {
			return
		}

		//
		var tagvals map[string]string
		newVal := reflect.NewAt(refval.Type(), unsafe.Pointer(refval.UnsafeAddr()))
		if tagvals = reflectutil.GetTagValues(ipakku.PAKKUTAG_AUTOWIRED, newVal); len(tagvals) == 0 {
			continue
		}

		//
		for field, valKey := range tagvals {
			var val interface{}
			if val, err = getModuleByName(valKey, l); nil != err {
				break
			}
			var ftype reflect.Type
			if ftype, err = reflectutil.GetStructFieldType(newVal, field); nil == err {
				logs.Infof("> Autowired %s <= %s[%s] \r\n", field, valKey, ftype.String())
				if ftype.Kind() != reflect.Interface {
					err = fmt.Errorf("only interface type injections are accepted, field: %s", field)
					break
				}
				if err = reflectutil.SetStructFieldValueUnSafe(newVal, field, val); nil != err {
					break
				}
			} else {
				logs.Infof("> Autowired %s <= %s[err=%s] \r\n", field, valKey, err.Error())
			}
		}

		// 再次递归
		if nil == err {
			err = autoWiredAnonymousStruct(newVal, l)
		}
	}
	return err
}

// getModuleByName 通过模块名字获取模块实例
func getModuleByName(valKey string, l ipakku.Loader) (val interface{}, err error) {
	if err = l.GetModuleByName(valKey, &val); nil != err && err.Error() == fmt.Sprintf(ipakku.ErrModuleNotFoundStr, valKey) {
		// 再从Params中找找看
		if val = l.GetParam(valKey).GetVal(); nil == val {
			err = fmt.Errorf(ipakku.ErrModuleNotFoundStr, valKey)
		} else {
			err = nil
		}
	}
	return
}
