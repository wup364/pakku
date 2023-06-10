// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 1. 通过重新实现 ixxx.go 接口
// 2. 在对应模块初始化之前注册实例 ipakku.Override.RegisterInterfaceImpl(val, interface-name, name) (如: init方法)
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

// Override 复写模块静态方法
var Override = overrideFuc{
	SetInterfaceDefaultImpl: doSetInterfaceDefaultImpl,
	RegisterInterfaceImpl:   doRegisterInterfaceImpl,
	AutowireInterfaceImpl:   doAutowireInterfaceImpl,
	GetImplementByName:      doGetImplementByName,
	SetModuleInfoImpl:       doSetModuleInfoImpl,
	GetModuleInfoImpl:       doGetModuleInfoImpl,
}

// overrideFuc 重载函数
type overrideFuc struct {
	GetImplementByName      func(interfaceName, implName, defaultImplName string) interface{}
	RegisterInterfaceImpl   func(val interface{}, interfaceName, implName string)
	AutowireInterfaceImpl   func(param paramGet, val interface{}, defaultName string) error
	SetInterfaceDefaultImpl func(param paramSet, interfaceName, implName string)
	SetModuleInfoImpl       func(val ModuleInfo)
	GetModuleInfoImpl       func() ModuleInfo
}

type paramGet interface {
	GetParam(key string) utypes.Object
}
type paramSet interface {
	SetParam(key string, val interface{})
}

// moduleInfoImpl moduleloader 默认的 ModuleInfo  记录器, 给他赋值以改变记录方式
var moduleInfoImpl ModuleInfo

// doSetModuleInfoImpl 注册模块信息记录实现方法
func doSetModuleInfoImpl(val ModuleInfo) {
	moduleInfoImpl = val
}

// doGetModuleInfoImpl 获取模块信息记录实现方法
func doGetModuleInfoImpl() ModuleInfo {
	return moduleInfoImpl
}

// implements 所有的ixxx.go实现实例, 结构: { ixxx: map[name]implement }
var implements = utypes.NewSafeMap()

// doGetImplementsByInterface 根据接口名字查找所有实现
func doGetImplementsByInterface(interfaceName string) *utypes.SafeMap {
	if val, ok := implements.Get(interfaceName); ok {
		return val.(*utypes.SafeMap)
	} else {
		newType := utypes.NewSafeMap()
		implements.Put(interfaceName, newType)
		return newType
	}
}

// doGetImplementByName 根据接口名字+(实例名字 || 默认实例名字), 获取具体实现对象
func doGetImplementByName(interfaceName, implName, defaultImplName string) interface{} {
	its := doGetImplementsByInterface(interfaceName)
	if val, ok := its.Get(implName); ok {
		return val
	} else if val, ok := its.Get(defaultImplName); ok {
		return val
	}
	return nil
}

// doRegisterInterfaceImpl 添加接口实现实例, interfaceName 接口名字需要和接口本身一致
func doRegisterInterfaceImpl(val interface{}, interfaceName, implName string) {
	doGetImplementsByInterface(interfaceName).Put(implName, val)
}

// doAutowireInterfaceImpl 多个相同接口下, 设置自动注入接口的实例名称
func doAutowireInterfaceImpl(param paramGet, val interface{}, defaultName string) error {
	var reft reflect.Type
	if reft = reflect.TypeOf(val); reft.Kind() != reflect.Ptr || reft.Elem().Kind() != reflect.Interface {
		return errors.New("only pointer interface are supported")
	}
	implName := moduleImplsPrefix + "." + reft.Elem().Name()
	impl := doGetImplementByName(reft.Elem().Name(), param.GetParam(implName).ToString(defaultName), defaultName)
	return reflectutil.SetInterfaceValueUnSafe(val, impl)
}

// doSetInterfaceDefaultImpl 设置默认接口实现, 在application实例上
func doSetInterfaceDefaultImpl(param paramSet, interfaceName, implName string) {
	param.SetParam(moduleImplsPrefix+"."+interfaceName, implName)
}
