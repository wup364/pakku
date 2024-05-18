// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of jsoncfg source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of jsoncfg software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and jsoncfg permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 入口包

package pakku

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/mloader"
	"github.com/wup364/pakku/modules/appcache"
	"github.com/wup364/pakku/modules/appconfig"
	"github.com/wup364/pakku/modules/appevent"
	"github.com/wup364/pakku/modules/appservice"
	"github.com/wup364/pakku/utils/fileutil"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/reflectutil"
)

// NewApplication 新建应用加载器
func NewApplication(name string) ipakku.ApplicationBootBuilder {
	boot := &ApplicationBootBuilder{
		locker:  new(sync.Mutex),
		modules: make([]ipakku.Module, 0),
		loader:  mloader.NewDefault(name),
	}
	//
	boot.mevent = &ModuleEventBuilder{boot: boot}
	boot.pkModules = &PakkuModuleBuilder{boot: boot}
	boot.csModules = &CustomModuleBuilder{boot: boot}
	boot.pkconf = &PakkuConfigureBuilder{
		boot:       boot,
		showBanner: true,
	}
	return boot
}

// ApplicationBootBuilder 程序启动引导
type ApplicationBootBuilder struct {
	locker    *sync.Mutex
	modules   []ipakku.Module
	loader    ipakku.Loader
	mevent    *ModuleEventBuilder
	pkModules *PakkuModuleBuilder
	csModules *CustomModuleBuilder
	pkconf    *PakkuConfigureBuilder
	pakapp    ipakku.PakkuApplication
}

// PakkuConfigure 应用配置操作
func (boot *ApplicationBootBuilder) PakkuConfigure() ipakku.PakkuConfigure {
	return boot.pkconf
}

// PakkuModules 默认模块启用操作
func (boot *ApplicationBootBuilder) PakkuModules() ipakku.PakkuModuleBuilder {
	return boot.pkModules
}

// CustomModules 自定义模块操作
func (boot *ApplicationBootBuilder) CustomModules() ipakku.CustomModuleBuilder {
	return boot.csModules
}

// ModuleEvents 模块事件监听器
func (boot *ApplicationBootBuilder) ModuleEvents() ipakku.ModuleEventBuilder {
	return boot.mevent
}

// Application 获取Application实例
func (boot *ApplicationBootBuilder) Application() ipakku.PakkuApplication {
	if boot.pakapp != nil {
		return boot.pakapp
	}

	return &PakkuApplication{
		Application: boot.loader.GetApplication(),
	}
}

// BootStart 启动程序&加载模块
func (boot *ApplicationBootBuilder) BootStart() ipakku.PakkuApplication {
	boot.locker.Lock()
	defer boot.locker.Unlock()

	if boot.pakapp == nil {
		if boot.pkconf.showBanner {
			boot.printBanner()
		}

		cwd, _ := os.Getwd()
		instanceID := strings.ToUpper(boot.loader.GetInstanceID())
		name := boot.loader.GetParam(ipakku.PARAMS_KEY_APPNAME).ToString(ipakku.DEFT_VAL_APPNAME)
		logs.Infof("New application, name: %s, pid: %d, cwd: %s, instance: %s \r\n", name, os.Getpid(), cwd, instanceID)
		boot.loader.Loads(boot.modules...)

		boot.pakapp = &PakkuApplication{
			Application: boot.loader.GetApplication(),
		}
	}

	return boot.pakapp
}

// addModule 加载模块
func (boot *ApplicationBootBuilder) addModule(mt ipakku.Module) {
	if !boot.modulesIsExist(mt) {
		boot.modules = append(boot.modules, mt)
	}
}

// addModules 加载模块
func (boot *ApplicationBootBuilder) addModules(mts ...ipakku.Module) {
	for i := 0; i < len(mts); i++ {
		if !boot.modulesIsExist(mts[i]) {
			boot.modules = append(boot.modules, mts[i])
		}
	}
}

// modulesIsExist 是否已经添加过了
func (boot *ApplicationBootBuilder) modulesIsExist(mt ipakku.Module) bool {
	if len(boot.modules) == 0 {
		return false
	}
	for i := 0; i < len(boot.modules); i++ {
		if boot.getModuleName(boot.modules[i]) == boot.getModuleName(mt) {
			return true
		}
	}
	return false
}

// getModuleName 获取模块名字(ID)
func (boot *ApplicationBootBuilder) getModuleName(mt ipakku.Module) string {
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

// printBanner 打印一些特殊记号
func (boot *ApplicationBootBuilder) printBanner() {
	bannerPath := "./.conf/banner.txt"
	if !fileutil.IsFile(bannerPath) {
		banner := "" +
			"              ,----------------,              ,---------, \r\n" +
			"         ,-----------------------,          ,\"        ,\"| \r\n" +
			"       ,\"                      ,\"|        ,\"        ,\"  | \r\n" +
			"      +-----------------------+  |      ,\"        ,\"    | \r\n" +
			"      |  .-----------------.  |  |     +---------+      | \r\n" +
			"      |  |                 |  |  |     | -==----'|      | \r\n" +
			"      |  |  I AM RUNNING!  |  |  |     |         |      | \r\n" +
			"      |  |  PLEASE WAIT .. |  |  |/----|'---=    |      | \r\n" +
			"      |  |  $ >_           |  |  |   ,/|==== ooo |      ; \r\n" +
			"      |  |                 |  |  |  // |(((( [33]|    ,\" \r\n" +
			"      |  '-----------------'  |,\" .;'| |((((     |  ,\" \r\n" +
			"      +-----------------------\"  ;;  | |         |,\" \r\n" +
			"         /_)______________(_/   /    | +---------+ \r\n" +
			"    ___________________________/___  ', \r\n" +
			"   /  oooooooooooooooo  .o.  oooo /,   \\,\"----------- \r\n" +
			"  / ==ooooooooooooooo==.o.  ooo= //   ,'\\--{)B     ,\" \r\n" +
			" /_==__==========__==_ooo__ooo=_/'   /___________,\" \r\n"
		if err := fileutil.WriteTextFile(bannerPath, banner); nil != err {
			logs.Errorln(err)
		} else {
			fmt.Println(banner)
		}
	} else {
		if banner, err := fileutil.ReadFileAsText(bannerPath); nil != err {
			logs.Errorln(err)
		} else {
			fmt.Println(banner)
		}
	}
}

// PakkuApplication 应用实例
type PakkuApplication struct {
	ipakku.Application
	pakgter ipakku.PakkuModulesGetter
}

// PakkuModules 获取默认携带的模块
func (pa *PakkuApplication) PakkuModules() ipakku.PakkuModulesGetter {
	if pa.pakgter == nil {
		pa.pakgter = &PakkuModulesGetter{app: pa.Application}
	}
	return pa.pakgter
}

// PakkuConfigureBuilder 应用配置
type PakkuConfigureBuilder struct {
	boot       *ApplicationBootBuilder
	showBanner bool
}

// SetLoggerOutput 设置日志输出方式
func (pkcf *PakkuConfigureBuilder) SetLoggerOutput(w io.Writer) ipakku.PakkuConfigure {
	logs.SetOutput(w)
	return pkcf
}

// SetLoggerLevel 设置日志输出级别 NONE DEBUG INFO ERROR
func (pkcf *PakkuConfigureBuilder) SetLoggerLevel(lv logs.LoggerLeve) ipakku.PakkuConfigure {
	logs.SetLoggerLevel(lv)
	return pkcf
}

// DisableBanner 禁止Banner输出
func (pkcf *PakkuConfigureBuilder) DisableBanner() ipakku.PakkuConfigure {
	pkcf.showBanner = false
	return pkcf
}

// PakkuModules 默认携带的模块
func (pkcf *PakkuConfigureBuilder) PakkuModules() ipakku.PakkuModuleBuilder {
	return pkcf.boot.pkModules
}

// PakkuModules 默认携带的模块
func (pkcf *PakkuConfigureBuilder) CustomModules() ipakku.CustomModuleBuilder {
	return pkcf.boot.csModules
}

// PakkuModuleBuilder 默认携带模块构造器
type PakkuModuleBuilder struct {
	boot *ApplicationBootBuilder
}

// EnableAppConfig 启用配置模块
func (pkm *PakkuModuleBuilder) EnableAppConfig() ipakku.PakkuModuleBuilder {
	pkm.boot.addModule(new(appconfig.AppConfig))
	return pkm
}

// EnableAppCache 启用缓存模块
func (pkm *PakkuModuleBuilder) EnableAppCache() ipakku.PakkuModuleBuilder {
	pkm.boot.addModule(new(appcache.AppCache))
	return pkm
}

// EnableAppEvent 启用事件模块
func (pkm *PakkuModuleBuilder) EnableAppEvent() ipakku.PakkuModuleBuilder {
	pkm.boot.addModule(new(appevent.AppEvent))
	return pkm
}

// EnableAppService 启用网络服务[WEB|RPC]模块
func (pkm *PakkuModuleBuilder) EnableAppService() ipakku.PakkuModuleBuilder {
	pkm.boot.addModule(new(appservice.AppService))
	return pkm
}

// PakkuModules 默认携带的模块
func (pkm *PakkuModuleBuilder) CustomModules() ipakku.CustomModuleBuilder {
	return pkm.boot.csModules
}

// ModuleEvents 模块事件监听器
func (pkm *PakkuModuleBuilder) ModuleEvents() ipakku.ModuleEventBuilder {
	return pkm.boot.mevent
}

// BootStart 启动程序&加载模块
func (pkm *PakkuModuleBuilder) BootStart() ipakku.PakkuApplication {
	return pkm.boot.BootStart()
}

// PakkuModulesGetter 获取默认携带的模块
type PakkuModulesGetter struct {
	app ipakku.Application
}

// GetAppConfig 获得配置模块
func (pg *PakkuModulesGetter) GetAppConfig() ipakku.AppConfig {
	var result ipakku.AppConfig
	if err := pg.app.Modules().GetModules(&result); nil != err {
		return nil
	}
	return result
}

// GetAppCache 获得缓存模块
func (pg *PakkuModulesGetter) GetAppCache() ipakku.AppCache {
	var result ipakku.AppCache
	if err := pg.app.Modules().GetModules(&result); nil != err {
		return nil
	}
	return result
}

// GetAppEvent 获得事件模块
func (pg *PakkuModulesGetter) GetAppEvent() ipakku.AppEvent {
	var result ipakku.AppEvent
	if err := pg.app.Modules().GetModules(&result); nil != err {
		return nil
	}
	return result
}

// GetAppService 获得网络服务[WEB|RPC]模块
func (pg *PakkuModulesGetter) GetAppService() ipakku.AppService {
	var result ipakku.AppService
	if err := pg.app.Modules().GetModules(&result); nil != err {
		return nil
	}
	return result
}

// CustomModuleBuilder 自定义模块构造器
type CustomModuleBuilder struct {
	boot *ApplicationBootBuilder
}

// AddModule 加载模块
func (cms *CustomModuleBuilder) AddModule(mt ipakku.Module) ipakku.CustomModuleBuilder {
	cms.boot.addModule(mt)
	return cms
}

// AddModules 加载模块
func (cms *CustomModuleBuilder) AddModules(mts ...ipakku.Module) ipakku.CustomModuleBuilder {
	cms.boot.addModules(mts...)
	return cms
}

// ModuleEvents 模块事件监听器
func (csm *CustomModuleBuilder) ModuleEvents() ipakku.ModuleEventBuilder {
	return csm.boot.mevent
}

// BootStart 启动程序&加载模块
func (csm *CustomModuleBuilder) BootStart() ipakku.PakkuApplication {
	return csm.boot.BootStart()
}

// ModuleEventBuilder 模块事件监听器
type ModuleEventBuilder struct {
	boot *ApplicationBootBuilder
}

// Listen 监听模块生命周期事件
func (meb ModuleEventBuilder) Listen(name string, event ipakku.ModuleEvent, val ipakku.OnModuleEvent) ipakku.ModuleEventBuilder {
	meb.boot.Application().Modules().OnModuleEvent(name, event, val)
	return meb
}

// BootStart 加载&启动程序
func (meb ModuleEventBuilder) BootStart() ipakku.PakkuApplication {
	return meb.boot.BootStart()
}
