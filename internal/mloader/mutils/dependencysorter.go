// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package mutils

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/reflectutil"
	"github.com/wup364/pakku/pkg/strutil"
)

// GetAutoWiredDependencies 获取依赖注入的依赖树信息
func GetAutoWiredDependencies(ptr ipakku.Module) (reuslt strutil.DS_M, err error) {
	// 仅支持指针类型结构体
	if t := reflect.TypeOf(ptr); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("only pointer object '%s' is supported", ipakku.STAG_AUTOWIRED)
		return
	}

	reuslt.Name = ptr.AsModule().Name
	var tagvals = make(map[string]string)
	if tagvals = reflectutil.GetTagValues(ipakku.STAG_AUTOWIRED, ptr); len(tagvals) == 0 {
		appendFormAnonymousStruct(ptr, &reuslt)
	} else {
		if err = appendDependencies(&reuslt, tagvals, ptr); nil != err {
			return
		}

		// 匿名嵌套结构体
		err = appendFormAnonymousStruct(ptr, &reuslt)
	}
	reuslt.Dependencies = strutil.RemoveDuplicatesAndEmpty(reuslt.Dependencies...)
	return
}

// appendFormAnonymousStruct 匿名嵌套结构体
func appendFormAnonymousStruct(ptr any, reuslt *strutil.DS_M) (err error) {
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

		// 递归
		var tagvals map[string]string
		newVal := reflect.NewAt(refval.Type(), unsafe.Pointer(refval.UnsafeAddr()))
		if tagvals = reflectutil.GetTagValues(ipakku.STAG_AUTOWIRED, newVal); len(tagvals) == 0 {
			if err = appendFormAnonymousStruct(newVal, reuslt); nil == err {
				continue
			}
			return
		}

		// append
		if err = appendDependencies(reuslt, tagvals, newVal); nil != err {
			return
		}

		// 递归
		err = appendFormAnonymousStruct(newVal, reuslt)
	}
	return err
}

// appendDependencies 追加依赖信息
func appendDependencies(reuslt *strutil.DS_M, tagvals map[string]string, ptr any) (err error) {
	if len(tagvals) == 0 {
		return
	}
	for field, valKey := range tagvals {
		if len(valKey) == 0 {
			var ftype reflect.Type
			if ftype, err = reflectutil.GetStructFieldType(ptr, field); nil != err {
				err = fmt.Errorf("autowire field %s is failed, error: %s", field, err.Error())
				break
			}
			reuslt.Dependencies = append(reuslt.Dependencies, ftype.Name())
		} else {
			reuslt.Dependencies = append(reuslt.Dependencies, valKey)
		}
	}

	return
}
