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
	"fmt"
	"reflect"
	"unsafe"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/reflectutil"
)

// AutoWired 自动注入依赖
func AutoWired(ptr interface{}, app ipakku.Application) (err error) {
	// 仅支持指针类型结构体
	if t := reflect.TypeOf(ptr); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("only pointer object '%s' is supported", ipakku.STAG_AUTOWIRED)
	}

	var tagvals = make(map[string]string)
	if tagvals = reflectutil.GetTagValues(ipakku.STAG_AUTOWIRED, ptr); len(tagvals) == 0 {
		return autoWiredAnonymousStruct(ptr, app)
	}

	// 执行自动注入
	if err = doAutowireFields(ptr, tagvals, app); nil == err {
		// 自动注入匿名嵌套结构体
		err = autoWiredAnonymousStruct(ptr, app)
	}

	return err
}

// autoWiredAnonymousStruct 自动注入匿名嵌套结构体
func autoWiredAnonymousStruct(ptr interface{}, app ipakku.Application) (err error) {
	var fields []reflect.StructField
	if fields = reflectutil.GetAnonymousOrNoneTypeNameField(ptr); len(fields) == 0 {
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
		if tagvals = reflectutil.GetTagValues(ipakku.STAG_AUTOWIRED, newVal); len(tagvals) == 0 {
			// 再次递归
			if err = autoWiredAnonymousStruct(newVal, app); nil == err {
				continue
			}
			return
		}

		// 执行自动注入
		if err = doAutowireFields(newVal, tagvals, app); nil != err {
			break
		}

		// 再次递归
		if err = autoWiredAnonymousStruct(newVal, app); nil != err {
			break
		}
	}
	return err
}

// doAutowireFields 执行对象下的字段自动注入
func doAutowireFields(ptr interface{}, tagvals map[string]string, app ipakku.Application) (err error) {
	if len(tagvals) == 0 {
		return
	}

	for field, valKey := range tagvals {
		var ftype reflect.Type
		if ftype, err = reflectutil.GetStructFieldType(ptr, field); nil != err {
			err = fmt.Errorf("autowire field %s is failed, error: %s", field, err.Error())
			break
		} else if ftype.Kind() != reflect.Interface {
			err = fmt.Errorf("autowire field %s is failed, error: only interface type injections are accepted", field)
			break
		}

		var val interface{}
		moduleName := getModuleName(valKey, ftype)
		if val, err = getModuleByName(moduleName, app); nil != err {
			break
		}
		if err = reflectutil.SetStructFieldValueUnSafe(ptr, field, val); nil != err {
			break
		}
		logs.Infof("> Autowired %s <= %s[%s] \r\n", field, moduleName, ftype.String())
	}
	return
}

// getModuleName 若name为空, 则取类型名字
func getModuleName(name string, ftype reflect.Type) string {
	if len(name) > 0 {
		return name
	}
	return ftype.Name()
}

// getModuleByName 通过模块名字获取模块实例
func getModuleByName(valKey string, app ipakku.Application) (val interface{}, err error) {
	if err = app.Modules().GetModuleByName(valKey, &val); nil != err && err.Error() == fmt.Sprintf(ipakku.ERR_MSG_MODULE_NOT_FOUND, valKey) {
		// 再从Params中找找看
		if val = app.Params().GetParam(valKey).GetVal(); nil == val {
			err = fmt.Errorf(ipakku.ERR_MSG_MODULE_NOT_FOUND, valKey)
		} else {
			err = nil
		}
	}
	return
}
