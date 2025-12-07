// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

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

	"github.com/wup364/pakku/internal/mloader/mutils"
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"
	"github.com/wup364/pakku/pkg/reflectutil"
	"github.com/wup364/pakku/pkg/strutil"
	"github.com/wup364/pakku/pkg/utypes"
)

// Loader 模块加载器, 实例化后可实现统一管理模板
type Loader struct {
	instanceID string                                          // 加载器实例ID
	events     *utypes.SafeMap[string, []ipakku.OnModuleEvent] // 模块生命周期事件
	modules    *utypes.SafeMap[string, ipakku.Module]          // 模块Map表
	mparams    *utypes.SafeMap[string, any]                    // 保存在模块对象中共享的字段key-value
	mrecord    ipakku.ModuleInfoRecorder                       // 模块信息记录器
}

// Loads 初始化模块(自动分析模块依赖), 初始化顺序: doReady -> doSetup -> doCheckVersion -> doInit -> doEnd
func (loader *Loader) Loads(mts ...ipakku.Module) {
	newMts := loader.dependencySort(mts...)
	for i := 0; i < len(newMts); i++ {
		loader.Load(newMts[i])
	}
}

// Load 初始化模块, 初始化顺序: doReady -> doSetup -> doCheckVersion -> doInit -> doEnd
func (loader *Loader) Load(mt ipakku.Module) {
	moduleOpts := mt.AsModule()
	moduleName := loader.getModuleName(mt)
	logs.Infof("> Loading %s Start ", moduleName)

	// doready 模块准备开始加载
	loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnReady)
	loader.doReady(moduleName, moduleOpts)

	// doSetup 模块安装
	var isSetup bool
	if isSetup = len(loader.GetModuleVersion(moduleName)) == 0; isSetup {
		loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnSetup)
		loader.doSetup(moduleName, moduleOpts)
		loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnSetupSucced)
	}

	// doCheckVersion 模块升级
	if !isSetup && loader.GetModuleVersion(moduleName) != strconv.FormatFloat(mt.AsModule().Version, 'f', 2, 64) {
		loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnUpdate)
		loader.doUpdate(moduleName, moduleOpts)
		loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnUpdateSucced)
	}

	// doInit 模块初始化
	loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnInit)
	loader.doInit(moduleName, moduleOpts)

	// doEnd 模块加载结束
	loader.modules.Put(moduleName, mt)
	logs.Infof("> Loading %s Complete ", moduleName)
	loader.doHandleModuleEvent(mt, ipakku.ModuleEventOnLoaded)
}

// Invoke 模块调用, 返回 []reflect.Value, 返回值暂时无法处理
func (loader *Loader) Invoke(name string, method string, params ...any) ([]reflect.Value, error) {
	if module, ok := loader.modules.Get(name); ok {
		val := reflect.ValueOf(module)
		fun := val.MethodByName(method)
		logs.Infof("> Invoke: "+name+"."+method+", %v, %+v ", fun, &fun)
		args := make([]reflect.Value, len(params))
		for i, temp := range params {
			args[i] = reflect.ValueOf(temp)
		}
		return fun.Call(args), nil
	}
	return nil, fmt.Errorf(ipakku.ERR_MSG_MODULE_NOT_FOUND, name)
}

// AutoWired 自动注入依赖对象
func (loader *Loader) AutoWired(structobj any) error {
	return mutils.AutoWired(structobj, loader)
}

// GetInstanceID 获取实例的ID
func (loader *Loader) GetInstanceID() string {
	return loader.instanceID
}

// SetParam 设置变量, 保存在模板加载器实例内部
func (loader *Loader) SetParam(key string, val any) {
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
		events = val
	} else {
		events = make([]ipakku.OnModuleEvent, 0)
	}
	loader.events.Put(eventKey, append(events, val))
}

// GetModuleByName 根据模块Name获取模块指针记录, 可以获取一个已经实例化的模块
func (loader *Loader) GetModuleByName(name string, val any) error {
	if tmp, ok := loader.modules.Get(name); ok {
		return reflectutil.SetInterfaceValueUnSafe(val, tmp)
	}
	return fmt.Errorf(ipakku.ERR_MSG_MODULE_NOT_FOUND, name)
}

// GetModules 获取模块, 模块名字和接口名字一样才能正常获得
func (loader *Loader) GetModules(val ...any) error {
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

// SetModuleInfoRecorder 设置模块信息记录器, 会自动调用init
func (loader *Loader) SetModuleInfoRecorder(mrecord ipakku.ModuleInfoRecorder) {
	if nil != mrecord {
		loader.mrecord = mrecord
		err := loader.mrecord.Init(loader.GetParam(ipakku.PARAMS_KEY_APPNAME).ToString(ipakku.DEFT_VAL_APPNAME))
		if nil != err {
			panic(err)
		}
	}
}

// GetModuleVersion 获取模块版本号
func (loader *Loader) GetModuleVersion(name string) string {
	return loader.mrecord.GetValue(name + ".SetupVer")
}

// GetApplication 获取当前实例
func (loader *Loader) GetApplication() ipakku.Application {
	return loader
}

// Params 保存实例中的键值对数据
func (loader *Loader) Params() ipakku.Params {
	return loader
}

// Modules 模块操作
func (loader *Loader) Modules() ipakku.Modules {
	return loader
}

// Utils 工具
func (loader *Loader) Utils() ipakku.Utils {
	return loader
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
		events = val
	}
	if val, ok := loader.events.Get(loader.getModuleEventKey(loader.getModuleName(mt), event)); ok {
		events = append(events, val...)
	}
	if len(events) > 0 {
		for i := 0; i < len(events); i++ {
			events[i](mt, loader)
		}
	}
}

// doReady 模块准备
func (loader *Loader) doReady(moduleId string, opts ipakku.Opts) {
	if nil != opts.OnReady {
		logs.Infof("> Execute %s.OnReady ", moduleId)
		opts.OnReady(loader)
	}
}

// doSetup 模块安装
func (loader *Loader) doSetup(moduleId string, opts ipakku.Opts) {
	if nil != opts.OnSetup {
		logs.Infof("> Execute %s.OnSetup ", moduleId)
		opts.OnSetup()
	}
	loader.setVersion(moduleId, opts.Version)
}

// doUpdate 模块升级
func (loader *Loader) doUpdate(moduleId string, opts ipakku.Opts) {
	if nil == opts.Updaters {
		loader.setVersion(moduleId, opts.Version)
		return
	}

	var updaters []ipakku.Updater
	if updaters = opts.Updaters(loader); len(updaters) == 0 {
		loader.setVersion(moduleId, opts.Version)
		return
	}

	var execList ipakku.Updaters = make([]ipakku.Updater, 0)
	if hv, err := strconv.ParseFloat(loader.GetModuleVersion(moduleId), 64); nil != err {
		logs.Panic(err)
	} else {
		for i := 0; i < len(updaters); i++ {
			if upv := updaters[i].Version(); upv > hv && upv <= opts.Version {
				execList = append(execList, updaters[i])
			}
		}
	}

	if len(execList) == 0 {
		loader.setVersion(moduleId, opts.Version)
		return
	}

	sort.Sort(execList)
	for i := 0; i < len(execList); i++ {
		logs.Infof("> Execute %s.Update ver=%.3f ", moduleId, execList[i].Version())
		if err := execList[i].Execute(loader); nil == err {
			loader.setVersion(moduleId, execList[i].Version())
			logs.Infof("> Completed %s.Update ver=%.3f ", moduleId, execList[i].Version())
		} else {
			logs.Panic(err)
		}
	}

	loader.setVersion(moduleId, opts.Version)
}

// doInit 模块初始化
func (loader *Loader) doInit(moduleId string, opts ipakku.Opts) {
	if nil != opts.OnInit {
		logs.Infof("> Execute %s.OnInit ", moduleId)
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

// dependencySort 依赖排序
func (loader *Loader) dependencySort(mts ...ipakku.Module) (res []ipakku.Module) {
	modules := []strutil.DS_M{}
	for i := 0; i < len(mts); i++ {
		if module, err := mutils.GetAutoWiredDependencies(mts[i]); nil != err {
			logs.Panic(err)
		} else {
			if len(module.Name) == 0 {
				module.Name = loader.getModuleName(mts[i])
			}
			modules = append(modules, module)
		}
	}

	sorted := strutil.DependencySorter(modules...)
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(mts); j++ {
			if loader.getModuleName(mts[j]) == sorted[i] {
				res = append(res, mts[j])
			}
		}
	}
	return
}
