// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 模块加载器
// 依赖包: utils.utypes.Object utils.strutil.strutil

package mloader

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/mloader/mutils"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/reflectutil"
	"github.com/wup364/pakku/utils/strutil"
	"github.com/wup364/pakku/utils/utypes"
)

// New 实例一个加载器对象
func New(name string) ipakku.Loader {
	loader := &Loader{
		events:     utypes.NewSafeMap(),
		modules:    utypes.NewSafeMap(),
		mparams:    utypes.NewSafeMap(),
		instanceID: strutil.GetUUID(),
	}
	if len(name) == 0 {
		name = ipakku.CONST_APPNAME
	}
	loader.SetParam(ipakku.PARAMKEY_APPNAME, name)
	loader.SetModuleInfoHandler(ipakku.Override.GetModuleInfoImpl())
	for _, v := range listeners {
		v.Bind(loader)
	}
	return loader
}

// Loader 模块加载器, 实例化后可实现统一管理模板
type Loader struct {
	instanceID string            // 加载器实例ID
	events     *utypes.SafeMap   // 模块生命周期事件
	modules    *utypes.SafeMap   // 模块Map表
	mparams    *utypes.SafeMap   // 保存在模块对象中共享的字段key-value
	mrecord    ipakku.ModuleInfo // 模块信息记录器
}

// Loads 初始化模块, 初始化顺序: doReady -> doSetup -> doCheckVersion -> doInit -> doEnd
func (loader *Loader) Loads(mts ...ipakku.Module) {
	for _, mt := range mts {
		loader.Load(mt)
	}
}

// Load 初始化模块, 初始化顺序: doReady -> doSetup -> doCheckVersion -> doInit -> doEnd
func (loader *Loader) Load(mt ipakku.Module) {
	moduleOpts := mt.AsModule()
	moduleName := loader.getModuleName(mt)
	logs.Infof("> Loading %s start \r\n", moduleName)

	// doready 模块准备开始加载
	loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnReady)
	loader.doReady(moduleName, moduleOpts)

	// doSetup 模块安装
	var isSetup bool
	if isSetup = len(loader.GetModuleVersion(moduleName)) == 0; isSetup {
		loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnSetup)
		loader.doSetup(moduleName, moduleOpts)
	}

	// doCheckVersion 模块升级
	if !isSetup && loader.GetModuleVersion(moduleName) != strconv.FormatFloat(mt.AsModule().Version, 'f', 2, 64) {
		loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnUpdate)
		loader.doUpdate(moduleName, moduleOpts)
	}

	// doInit 模块初始化
	loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnInit)
	loader.doInit(moduleName, moduleOpts)

	// doEnd 模块加载结束
	loader.modules.Put(moduleName, mt)
	logs.Infof("> Loading %s complete \r\n", moduleName)
	loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnLoaded)
}

// Invoke 模块调用, 返回 []reflect.Value, 返回值暂时无法处理
func (loader *Loader) Invoke(name string, method string, params ...interface{}) ([]reflect.Value, error) {
	if module, ok := loader.modules.Get(name); ok {
		val := reflect.ValueOf(module)
		fun := val.MethodByName(method)
		logs.Infof("> Invoke: "+name+"."+method+", %v, %+v \r\n", fun, &fun)
		args := make([]reflect.Value, len(params))
		for i, temp := range params {
			args[i] = reflect.ValueOf(temp)
		}
		return fun.Call(args), nil
	}
	return nil, fmt.Errorf(ipakku.ErrModuleNotFoundStr, name)
}

// AutoWired 自动注入依赖对象
func (loader *Loader) AutoWired(structobj interface{}) error {
	return mutils.AutoWired(structobj, loader)
}

// GetInstanceID 获取实例的ID
func (loader *Loader) GetInstanceID() string {
	return loader.instanceID
}

// SetParam 设置变量, 保存在模板加载器实例内部
func (loader *Loader) SetParam(key string, val interface{}) {
	loader.mparams.Put(key, val)
}

// GetParam 模板加载器实例上的变量
func (loader *Loader) GetParam(key string) utypes.Object {
	if val, ok := loader.mparams.Get(key); ok {
		return utypes.NewObject(val)
	}
	return utypes.Object{}
}

// OnModuleEvent 监听模块生命周期事件
func (loader *Loader) OnModuleEvent(name string, event ipakku.ModuleEvent, val ipakku.OnModuleEvent) {
	var events []ipakku.OnModuleEvent
	eventKey := loader.getModuleEventKey(name, event)
	if val, ok := loader.events.Get(eventKey); ok {
		events = val.([]ipakku.OnModuleEvent)
	} else {
		events = make([]ipakku.OnModuleEvent, 0)
	}
	loader.events.Put(eventKey, append(events, val))
}

// GetModuleByName 根据模块Name获取模块指针记录, 可以获取一个已经实例化的模块
func (loader *Loader) GetModuleByName(name string, val interface{}) error {
	if tmp, ok := loader.modules.Get(name); ok {
		return reflectutil.SetInterfaceValueUnSafe(val, tmp)
	}
	return fmt.Errorf(ipakku.ErrModuleNotFoundStr, name)
}

// GetModules 获取模块, 模块名字和接口名字一样才能正常获得
func (loader *Loader) GetModules(val ...interface{}) error {
	for i := 0; i < len(val); i++ {
		if nil == val[i] {
			return errors.New("the input object must be pointer interface, can not be nil")

		} else if valType := reflect.TypeOf(val[i]); valType.Kind() != reflect.Ptr || valType.Elem().Kind() != reflect.Interface {
			return errors.New("the input object must be pointer interface")

		} else if err := loader.GetModuleByName(valType.Elem().Name(), val[i]); nil != err {
			return err
		}
	}
	return nil
}

// SetModuleInfoHandler 设置模块信息记录器, 会自动调用init
func (loader *Loader) SetModuleInfoHandler(mrecord ipakku.ModuleInfo) {
	if nil != mrecord {
		loader.mrecord = mrecord
		err := loader.mrecord.Init(loader.GetParam(ipakku.PARAMKEY_APPNAME).ToString("app"))
		if nil != err {
			panic(err)
		}
	}
}

// GetModuleVersion 获取模块版本号
func (loader *Loader) GetModuleVersion(name string) string {
	return loader.mrecord.GetValue(name + ".SetupVer")
}

// setVersion 设置模块版本号 - 模块保留小数两位
func (loader *Loader) setVersion(moduleName string, version float64) {
	loader.mrecord.SetValue(moduleName+".SetupVer", strconv.FormatFloat(version, 'f', 2, 64))
	loader.mrecord.SetValue(moduleName+".SetupDate", strconv.FormatInt(time.Now().UnixNano(), 10))
}

// doHandleModuleEvent 执行监听模块生命周期事件
func (loader *Loader) doHandleModuleEvent(mt ipakku.Module, event ipakku.ModuleEvent) {
	var events []ipakku.OnModuleEvent
	if val, ok := loader.events.Get(loader.getModuleEventKey("*", event)); ok {
		if funs, ok := val.([]ipakku.OnModuleEvent); ok {
			events = funs
		}
	}
	if val, ok := loader.events.Get(loader.getModuleEventKey(loader.getModuleName(mt), event)); ok {
		if funs, ok := val.([]ipakku.OnModuleEvent); ok {
			events = append(events, funs...)
		}
	}
	if len(events) > 0 {
		for i := 0; i < len(events); i++ {
			events[i](mt, loader)
		}
	}
}

// doReady 模块准备
func (loader *Loader) doReady(moduleName string, opts ipakku.Opts) {
	if nil != opts.OnReady {
		logs.Infof("> Execute Module.OnReady \r\n")
		opts.OnReady(loader)
	}
}

// doSetup 模块安装
func (loader *Loader) doSetup(moduleName string, opts ipakku.Opts) {
	if nil != opts.OnSetup {
		logs.Infof("> Execute Module.OnSetup \r\n")
		opts.OnSetup()
	}
	loader.setVersion(moduleName, opts.Version)
}

// doUpdate 模块升级
func (loader *Loader) doUpdate(moduleName string, opts ipakku.Opts) {
	if nil == opts.Updaters {
		loader.setVersion(moduleName, opts.Version)
		return
	}

	var updaters []ipakku.Updater
	if updaters = opts.Updaters(loader); len(updaters) == 0 {
		loader.setVersion(moduleName, opts.Version)
		return
	}

	var execList ipakku.Updaters = make([]ipakku.Updater, 0)
	if hv, err := strconv.ParseFloat(loader.GetModuleVersion(moduleName), 64); nil != err {
		logs.Panicln(err)
	} else {
		for i := 0; i < len(updaters); i++ {
			if upv := updaters[i].Version(); upv > hv && upv <= opts.Version {
				execList = append(execList, updaters[i])
			}
		}
	}

	if len(execList) == 0 {
		loader.setVersion(moduleName, opts.Version)
		return
	}

	sort.Sort(execList)
	for i := 0; i < len(execList); i++ {
		logs.Infof("> Execute Module.Update ver=%.3f \r\n", execList[i].Version())
		if err := execList[i].Execute(loader); nil == err {
			loader.setVersion(moduleName, execList[i].Version())
			logs.Infof("> Completed Module.Update ver=%.3f \r\n", execList[i].Version())
		} else {
			logs.Panicln(err)
		}
	}

	loader.setVersion(moduleName, opts.Version)
}

// doInit 模块初始化
func (loader *Loader) doInit(moduleName string, opts ipakku.Opts) {
	if nil != opts.OnInit {
		logs.Infof("> Execute Module.OnInit \r\n")
		opts.OnInit()
	}
}

// getModuleEventKey getModuleEventKey
func (loader *Loader) getModuleEventKey(name string, event ipakku.ModuleEvent) string {
	return "ModuleEvent." + name + "." + string(event)
}

// getModuleName 获取模块名字(ID)
func (loader *Loader) getModuleName(mt ipakku.Module) string {
	if moduleName := mt.AsModule().Name; len(moduleName) == 0 {
		if mtype := reflectutil.GetNotPtrRefType(mt); nil == mtype {
			panic(fmt.Errorf("unable to obtain this object type: %T", mt))
		} else {
			return mtype.Name()
		}
	} else {
		return moduleName
	}
}
