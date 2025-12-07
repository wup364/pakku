// SPDX-License-Identifier: MIT
// Copyright (C) 2025 WuPeng <wup364@outlook.com>.

// 静态页面加载器
package appstaticpage

import (
	"strings"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"
	"github.com/wup364/pakku/pkg/strutil"

	"net/http"
)

// StaticPageLoader 静态资源加载器
type StaticPageLoader struct {
	sv ipakku.AppService `@autowired:""`
}

// AsModule 模块加载器接口实现, 返回模块信息&配置
func (staticPage *StaticPageLoader) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Version:     1.0,
		Description: "StaticPage module",
		OnReady: func(app ipakku.Application) {
			if pageConfig, err := GetStaticPageConfig(); nil != err {
				logs.Panicf("Failed to load staticPage configuration: %v", err)
			} else if nil != pageConfig {
				staticPage.registerStaticPages(*pageConfig)
			}
		},
	}
}

// registerStaticPages 注册静态页面配置
func (staticPage *StaticPageLoader) registerStaticPages(config StaticPageConfig) {
	// 1. 处理重定向
	if len(config.Redirect.Path) > 0 && len(config.Redirect.Target) > 0 && config.Redirect.Status > 0 {
		logs.Debugf("Register static page redirect: %s -> %s (status: %d)", config.Redirect.Path, config.Redirect.Target, config.Redirect.Status)
		if err := staticPage.sv.Get(config.Redirect.Path, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, config.Redirect.Target, config.Redirect.Status)
		}); nil != err {
			logs.Panic(err)
		}
	}

	// 2. 注册静态文件
	for _, file := range config.StaticFiles {
		logs.Debugf("Register static page file: %s -> %s", file.Path, file.FilePath)
		if err := staticPage.sv.SetStaticFile(file.Path, file.FilePath, nil); nil != err {
			logs.Panic(err)
		}
	}

	// 3. 注册静态目录
	if len(config.StaticDirectory.Path) > 0 && len(config.StaticDirectory.Directory) > 0 {
		logs.Debugf("Register static page directory: %s -> %s", config.StaticDirectory.Path, config.StaticDirectory.Directory)
		if err := staticPage.sv.SetStaticDIR(config.StaticDirectory.Path, config.StaticDirectory.Directory, nil); nil != err {
			logs.Panic(err)
		}
	}

	// 4. 注册过滤器
	if len(config.StaticDirectory.Path) > 0 {
		logs.Debugf("Register static page filter: %s", config.StaticDirectory.Path)
		config.StaticDirectory.Path = strutil.Parse2UnixPath(config.StaticDirectory.Path)
		if err := staticPage.sv.Filter(config.StaticDirectory.Path, func(w http.ResponseWriter, r *http.Request) bool {
			return staticPage.staticFilter(w, r, config)
		}); nil != err {
			logs.Panic(err)
		}
		if err := staticPage.sv.Filter(config.StaticDirectory.Path+"/:**", func(w http.ResponseWriter, r *http.Request) bool {
			return staticPage.staticFilter(w, r, config)
		}); nil != err {
			logs.Panic(err)
		}
	}
}

// staticFilter 静态资源过滤器
func (staticPage *StaticPageLoader) staticFilter(w http.ResponseWriter, r *http.Request, config StaticPageConfig) bool {

	// OPT 跨域预检
	if config.EnableCORS && r.Method == http.MethodOptions {
		if len(config.AllowedOrigins) > 0 {
			origin := r.Header.Get("Origin")
			if strutil.EqualsAny(origin, config.AllowedOrigins...) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		allowedHeaders := config.AllowedHeaders
		if len(allowedHeaders) == 0 {
			allowedHeaders = r.Header["Access-Control-Request-Headers"]
		}
		if len(allowedHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
		}

		allowedMethods := config.AllowedMethods
		if len(allowedMethods) == 0 {
			allowedMethods = r.Header["Access-Control-Request-Method"]
		}
		if len(allowedMethods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
		}

		w.WriteHeader(http.StatusOK)
		return false
	}

	// 检查请求方法是否允许
	if len(config.AllowedMethods) > 0 && !strutil.EqualsAny(r.Method, config.AllowedMethods...) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}

	return true
}
