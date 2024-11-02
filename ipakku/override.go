// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 1. 通过重新实现 ixxx.go 接口
// 2. 在对应模块初始化之前注册实例 ipakku.Override.RegisterPakkuModuleImplement(val, interface-name, name) (如: init方法)
// 3. 再在启动时app.SetParam(key, name)就可以替代默认模块啦~
package ipakku

import (
	"errors"
	"reflect"

	"github.com/wup364/pakku/utils/reflectutil"
	"github.com/wup364/pakku/utils/utypes"
)

const (
	// moduleImplsPrefix 模块接口实现注册前缀, 如: pakku.module.implement.IConfig
	moduleImplsPrefix = "pakku.module.implement"
)

// PakkuConf  PakkuConf配置, 如: 复写模块、设置模块信息记录器
var PakkuConf = pakkuConfFuc{
	// SetPakkuModuleImplement 设置默认接口实现, 在application实例上
	SetPakkuModuleImplement: doSetPakkuModuleImplement,

	// GetPakkuModuleImplement 根据接口名字+(实例名字 || 默认实例名字), 获取具体实现对象
	GetPakkuModuleImplement: doGetPakkuModuleImplement,

	// RegisterPakkuModuleImplement 添加接口实现实例, interfaceName 接口名字需要和接口本身一致
	RegisterPakkuModuleImplement: doRegisterPakkuModuleImplement,

	// AutowirePakkuModuleImplement 多个相同接口下, 设置自动注入接口的实例名称
	AutowirePakkuModuleImplement: doAutowirePakkuModuleImplement,

	// SetModuleInfoRecorderImplement 模块信息记录实现方法
	SetModuleInfoRecorderImplement: doSetModuleInfoRecorderImplement,

	// GetModuleInfoRecorderImplement 获取模块信息记录实现方法
	GetModuleInfoRecorderImplement: doGetModuleInfoRecorderImplement,
}

// moduleInfoImpl moduleloader 默认的 ModuleInfo  记录器, 给他赋值以改变记录方式
var moduleInfoImpl ModuleInfoRecorder

// implements 所有的ixxx.go实现实例, 结构: { ixxx: map[name]implement }
var implements = utypes.NewSafeMap[string, *utypes.SafeMap[string, any]]()

// pakkuConfFuc 重载函数
type pakkuConfFuc struct {
	// SetPakkuModuleImplement 设置默认接口实现, 在application实例上
	SetPakkuModuleImplement func(param ParamSetter, interfaceName, name string)

	// GetPakkuModuleImplement 根据接口名字+(实例名字 || 默认实例名字), 获取具体实现对象
	GetPakkuModuleImplement func(interfaceName, name, defaultName string) any

	// RegisterPakkuModuleImplement 添加接口实现实例, interfaceName 接口名字需要和接口本身一致
	RegisterPakkuModuleImplement func(val any, interfaceName, name string)

	// AutowirePakkuModuleImplement 多个相同接口下, 设置自动注入接口的实例名称
	AutowirePakkuModuleImplement func(param ParamGetter, name any, defaultName string) error

	// SetModuleInfoRecorderImplement 模块信息记录实现方法
	SetModuleInfoRecorderImplement func(val ModuleInfoRecorder)

	// GetModuleInfoRecorderImplement 获取模块信息记录实现方法
	GetModuleInfoRecorderImplement func() ModuleInfoRecorder
}

// doSetModuleInfoRecorderImplement 注册模块信息记录实现方法
func doSetModuleInfoRecorderImplement(val ModuleInfoRecorder) {
	moduleInfoImpl = val
}

// doGetModuleInfoRecorderImplement 获取模块信息记录实现方法
func doGetModuleInfoRecorderImplement() ModuleInfoRecorder {
	return moduleInfoImpl
}

// doGetImplementsByInterface 根据接口名字查找所有实现
func doGetImplementsByInterface(interfaceName string) *utypes.SafeMap[string, any] {
	if val, ok := implements.Get(interfaceName); ok {
		return val
	} else {
		newType := utypes.NewSafeMap[string, any]()
		implements.Put(interfaceName, newType)
		return newType
	}
}

// doGetPakkuModuleImplement 根据接口名字+(实例名字 || 默认实例名字), 获取具体实现对象
func doGetPakkuModuleImplement(interfaceName, implName, defaultImplName string) any {
	its := doGetImplementsByInterface(interfaceName)
	if val, ok := its.Get(implName); ok {
		return val
	} else if val, ok := its.Get(defaultImplName); ok {
		return val
	}
	return nil
}

// doRegisterPakkuModuleImplement 添加接口实现实例, interfaceName 接口名字需要和接口本身一致
func doRegisterPakkuModuleImplement(val any, interfaceName, name string) {
	doGetImplementsByInterface(interfaceName).Put(name, val)
}

// doAutowirePakkuModuleImplement 多个相同接口下, 设置自动注入接口的实例名称
func doAutowirePakkuModuleImplement(param ParamGetter, name any, defaultName string) error {
	var reft reflect.Type
	if reft = reflect.TypeOf(name); reft.Kind() != reflect.Ptr || reft.Elem().Kind() != reflect.Interface {
		return errors.New("only pointer interface are supported")
	}
	implName := moduleImplsPrefix + "." + reft.Elem().Name()
	impl := doGetPakkuModuleImplement(reft.Elem().Name(), param.GetParam(implName).ToString(defaultName), defaultName)
	return reflectutil.SetInterfaceValueUnSafe(name, impl)
}

// doSetPakkuModuleImplement 设置默认接口实现, 在application实例上
func doSetPakkuModuleImplement(param ParamSetter, interfaceName, name string) {
	param.SetParam(moduleImplsPrefix+"."+interfaceName, name)
}
