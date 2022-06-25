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

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/mloader"
	"github.com/wup364/pakku/modules/appcache"
	"github.com/wup364/pakku/modules/appconfig"
	"github.com/wup364/pakku/modules/appevent"
	"github.com/wup364/pakku/modules/appservice"
	"github.com/wup364/pakku/utils/fileutil"
	"github.com/wup364/pakku/utils/logs"
)

// NewApplication 新建应用加载器
func NewApplication(name string) ipakku.ApplicationBoot {
	return &ApplicationBoot{
		modules:    make([]ipakku.Module, 0),
		showBanner: true,
		application: &Application{
			Loader: mloader.New(name),
		},
	}
}

// ApplicationBoot 程序启动引导
type ApplicationBoot struct {
	modules     []ipakku.Module
	application *Application
	showBanner  bool
	loadedCore  bool
	loadedNet   bool
}

// GetInstanceID 获取实例的ID
func (boot *ApplicationBoot) GetApplication() ipakku.Application {
	return boot.application
}

// AddModule 加载模块
func (boot *ApplicationBoot) AddModule(mt ipakku.Module) ipakku.ApplicationBoot {
	boot.modules = append(boot.modules, mt)
	return boot
}

// AddModules 加载模块
func (boot *ApplicationBoot) AddModules(mts ...ipakku.Module) ipakku.ApplicationBoot {
	boot.modules = append(boot.modules, mts...)
	return boot
}

// EnableCoreModule 加载基础的运行模块 [AppConfig, AppCache, AppEvent]
func (boot *ApplicationBoot) EnableCoreModule() ipakku.ApplicationBoot {
	if !boot.loadedCore {
		boot.loadedCore = true
		boot.AddModules(new(appconfig.AppConfig), new(appcache.AppCache), new(appevent.AppEvent))
	}
	return boot
}

// EnableNetModule 启用网络模块 [AppService]
func (boot *ApplicationBoot) EnableNetModule() ipakku.ApplicationBoot {
	if !boot.loadedNet {
		boot.loadedNet = true
		boot.AddModule(&appservice.AppService{})
	}
	return boot
}

// SetLoggerOutput 设置日志输出方式
func (boot *ApplicationBoot) SetLoggerOutput(w io.Writer) ipakku.ApplicationBoot {
	logs.SetOutput(w)
	return boot
}

// SetLoggerLevel 设置日志输出级别 NONE DEBUG INFO ERROR
func (boot *ApplicationBoot) SetLoggerLevel(lv logs.LoggerLeve) ipakku.ApplicationBoot {
	logs.SetLoggerLevel(lv)
	return boot
}

// DisableBanner 禁止Banner输出
func (boot *ApplicationBoot) DisableBanner() ipakku.ApplicationBoot {
	boot.showBanner = false
	return boot
}

// BootStart 启动程序&加载模块
func (boot *ApplicationBoot) BootStart() ipakku.Application {
	if boot.showBanner {
		boot.printBanner()
	}
	cwd, _ := os.Getwd()
	app := boot.GetApplication()
	name := app.GetParam(ipakku.PARAMKEY_APPNAME).ToString("")
	logs.Infof("New application, pid: %d, name: %s, cwd: %s, instance: %s \r\n", os.Getpid(), name, cwd, strings.ToUpper(app.GetInstanceID()))
	return app.LoadModules(boot.modules...)
}

// printBanner 打印一些特殊记号
func (app *ApplicationBoot) printBanner() {
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

// Application 暴露loader对象的core方法
type Application struct {
	ipakku.Loader
}

// LoadModule 加载模块
func (app *Application) LoadModule(mt ipakku.Module) ipakku.Application {
	app.Loader.Load(mt)
	return app
}

// LoadModules 加载模块
func (app *Application) LoadModules(mts ...ipakku.Module) ipakku.Application {
	app.Loader.Loads(mts...)
	return app
}
