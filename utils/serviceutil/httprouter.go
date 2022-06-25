// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package serviceutil

// http服务器工具-URL路由管理
// 请求处理逻辑: 过滤器 > 路径匹配 > END
// 过滤器优先级: 全匹配url > 正则url > 全局设定 > 无匹配(next)
// 路径处理优先级: 全匹配url > 正则url > 默认设定 > 无匹配(404)
// isDebug 参数在生产环境注意关闭, 有成倍性能差距 6000/sec-> 25000/sec

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/strutil"
	"github.com/wup364/pakku/utils/utypes"
)

// HandlerFunc 定义请求处理器
type HandlerFunc func(http.ResponseWriter, *http.Request)

// FilterFunc http请求过滤器, 返回bool, true: 继续, false: 停止
type FilterFunc func(http.ResponseWriter, *http.Request) bool

// RuntimeErrorHandlerFunc 未知异常处理函数
type RuntimeErrorHandlerFunc func(http.ResponseWriter, *http.Request, interface{})

// ServiceRouter 实现了http.Server接口的ServeHTTP方法
type ServiceRouter struct {
	initial             bool                    // 是否初始化过了
	isDebug             bool                    // 调试模式可以打印信息
	runtimeError        RuntimeErrorHandlerFunc // 未知异常处理函数
	defaultHandler      HandlerFunc             // 默认的url处理, 可以用于处理静态资源
	urlHandlers         *utypes.SafeMap         // url路径全匹配路由表
	regexpHandlers      *utypes.SafeMap         // url路径正则配路由表
	regexpHandlersIndex []string                // url路径正则配路由表-索引(用于保存顺序)
	regexpFilters       *utypes.SafeMap         // url路径正则匹配过滤器
	regexpFiltersIndex  []string                // url路径正则匹配过滤器-索引(用于保存顺序)
	defaultFileter      FilterFunc
}

// ServiceRouter 根据注册的路由表调用对应的函数
// 优先匹配全url > 正则url > 默认处理器 > 404
func (serviceRouter *ServiceRouter) doHandle(w http.ResponseWriter, r *http.Request) {
	surl, err := buildHandlerURL(r.Method, r.URL.Path)
	if nil != err {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 1.0 如果是url全匹配, 则直接执行handler函数 - 有指定请求方式, POST, GET..
	if h, ok := serviceRouter.urlHandlers.Get(surl); ok {
		if serviceRouter.isDebug {
			logs.Debugln("[URL.Handler.Path]", surl)
		}
		h.(HandlerFunc)(w, r)
		return
	}
	// 1.1 如果是url全匹配, 则直接执行handler函数 - 无指定请求方式, POST, GET..
	anyurl, _ := buildHandlerURL("", r.URL.Path)
	if h, ok := serviceRouter.urlHandlers.Get(anyurl); ok {
		if serviceRouter.isDebug {
			logs.Debugln("[URL.Handler.Path]", anyurl)
		}
		h.(HandlerFunc)(w, r)
		return
	}
	// 2.0 如果是url正则检查, 则需要检查正则, 正则为':'后面的字符
	for i := 0; i < len(serviceRouter.regexpHandlersIndex); i++ {
		symbolIndex := strings.Index(serviceRouter.regexpHandlersIndex[i], ":")
		if symbolIndex == -1 {
			continue
		}
		matched := false
		baseURL := serviceRouter.regexpHandlersIndex[i][:symbolIndex]
		if strings.HasPrefix(surl, baseURL) {
			matched, _ = regexp.MatchString(baseURL+serviceRouter.regexpHandlersIndex[i][symbolIndex+1:], surl)
		} else if strings.HasPrefix(anyurl, baseURL) {
			matched, _ = regexp.MatchString(baseURL+serviceRouter.regexpHandlersIndex[i][symbolIndex+1:], anyurl)
		}
		if matched {
			if handler, ok := serviceRouter.regexpHandlers.Get(serviceRouter.regexpHandlersIndex[i]); ok {
				if serviceRouter.isDebug {
					logs.Debugln("[URL.Handler.Regexp]", surl, serviceRouter.regexpHandlersIndex[i])
				}
				handler.(HandlerFunc)(w, r)
				return
			}
		}
	}
	// 没有注册的地址, 使用默认处理器
	if serviceRouter.defaultHandler != nil {
		if serviceRouter.isDebug {
			logs.Debugln("[URL.Handler.Default]", surl)
		}
		serviceRouter.defaultHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// doFilter 根据注册的过滤器表调用对应的函数
// 优先匹配全url > 正则url > 全局过滤器 > 直接通过
func (serviceRouter *ServiceRouter) doFilter(w http.ResponseWriter, r *http.Request) {
	// 1 是否有指定路径的正则匹配过滤器设定, 取出所有
	if nil != serviceRouter.regexpFiltersIndex && len(serviceRouter.regexpFiltersIndex) > 0 {
		matched := make([]string, 0)
		for i := 0; i < len(serviceRouter.regexpFiltersIndex); i++ {
			symbolIndex := strings.Index(serviceRouter.regexpFiltersIndex[i], ":")
			if symbolIndex == -1 {
				if r.URL.Path == serviceRouter.regexpFiltersIndex[i] {
					matched = append(matched, serviceRouter.regexpFiltersIndex[i])
				}
				continue
			}
			baseURL := serviceRouter.regexpFiltersIndex[i][:symbolIndex]
			if !strings.HasPrefix(r.URL.Path, baseURL) {
				continue
			}
			if ok, _ := regexp.MatchString(serviceRouter.regexpFiltersIndex[i][:symbolIndex]+serviceRouter.regexpFiltersIndex[i][symbolIndex+1:], r.URL.Path); ok {
				matched = append(matched, serviceRouter.regexpFiltersIndex[i])
			}
		}
		if len(matched) > 0 {
			for i := 0; i < len(matched); i++ {
				if h, ok := serviceRouter.regexpFilters.Get(matched[i]); ok {
					if serviceRouter.isDebug {
						logs.Debugln("[URL.Filter]", matched[i])
					}
					if !h.(FilterFunc)(w, r) {
						return
					}
				}
			}
			// 符合所有过滤器要求
			serviceRouter.doHandle(w, r)
			return
		}
	}
	// 2. 检擦是否有全局过滤器存在, 如果有则执行它
	if nil != serviceRouter.defaultFileter {
		if serviceRouter.isDebug {
			logs.Debugln("[URL.Filter.Default]", r.URL.Path)
		}
		if serviceRouter.defaultFileter(w, r) {
			serviceRouter.doHandle(w, r)
		}
		return
	}
	// 3. 啥也没有设定
	serviceRouter.doHandle(w, r)
}

// ServeHTTP 实现http.Server接口的ServeHTTP方法
func (serviceRouter *ServiceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if serviceRouter.runtimeError != nil {
			if err := recover(); nil != err {
				serviceRouter.runtimeError(w, r, err)
			}
		}
	}()
	serviceRouter.checkInit()
	// 处理前进行过滤处理
	serviceRouter.doFilter(w, r)
}

// ClearHandlersMap 清空路由表
func (serviceRouter *ServiceRouter) ClearHandlersMap() {
	serviceRouter.urlHandlers.Clear()
	serviceRouter.regexpHandlers.Clear()
	serviceRouter.regexpHandlersIndex = make([]string, 0)
}

// SetDebug 是否输出url请求信息
func (serviceRouter *ServiceRouter) SetDebug(isDebug bool) {
	serviceRouter.isDebug = isDebug
}

// SetDefaultHandler 设置默认相应函数, 当无匹配时触发
func (serviceRouter *ServiceRouter) SetDefaultHandler(defaultHandler HandlerFunc) {
	if serviceRouter.isDebug {
		logs.Debugln("Set default handler")
	}
	serviceRouter.defaultHandler = defaultHandler
}

// SetRuntimeErrorHandler 设置全局handler执行异常捕获
func (serviceRouter *ServiceRouter) SetRuntimeErrorHandler(h RuntimeErrorHandlerFunc) {
	serviceRouter.runtimeError = h
}

// SetDefaultFilter 设置默认过滤器, 设置后, 如果不调用next函数则不进行下一步处理
// type FilterFunc func(http.ResponseWriter, *http.Request, func( ))
func (serviceRouter *ServiceRouter) SetDefaultFilter(globalFilter FilterFunc) {
	if serviceRouter.isDebug {
		logs.Debugln("Set default filter")
	}
	serviceRouter.defaultFileter = globalFilter
}

// AddURLFilter 设置url过滤器, 设置后, 如果不调用next函数则不进行下一步处理
// 过滤器有优先调用权, 正则匹配路径有先后顺序
// type FilterFunc func(http.ResponseWriter, *http.Request, func( ))
func (serviceRouter *ServiceRouter) AddURLFilter(url string, filter FilterFunc) error {
	if len(url) == 0 {
		return errors.New("filter url is empty")
	}
	serviceRouter.checkInit()
	url = strutil.Parse2UnixPath(url)
	if serviceRouter.isDebug {
		logs.Debugln("AddURLFilter: ", url)
	}
	serviceRouter.regexpFilters.Put(url, filter)
	serviceRouter.addFilterIndex(url)
	return nil
}

// addFilterIndex 添加filter索引, 安装路径长度排序
func (serviceRouter *ServiceRouter) addFilterIndex(url string) {
	newArray := append(serviceRouter.regexpFiltersIndex, url)
	strutil.SortByLen(newArray, true)
	strutil.SortBySplitLen(newArray, "/", true)
	serviceRouter.regexpFiltersIndex = newArray
}

// removeFilterIndex 删除filter索引
func (serviceRouter *ServiceRouter) removeFilterIndex(url string) {
	if len(url) > 0 {
		for i := 0; i < len(serviceRouter.regexpFiltersIndex); i++ {
			if serviceRouter.regexpFiltersIndex[i] == url {
				serviceRouter.regexpFiltersIndex = append(serviceRouter.regexpFiltersIndex[:i], serviceRouter.regexpFiltersIndex[i+i:]...)
				break
			}
		}
	}
}

// RemoveFilter 删除一个过滤器
func (serviceRouter *ServiceRouter) RemoveFilter(url string) {
	if len(url) == 0 {
		return
	}
	if serviceRouter.isDebug {
		logs.Debugln("RemoveFilter:", url)
	}
	if nil != serviceRouter.regexpFilters {
		if serviceRouter.regexpFilters.ContainsKey(url) {
			serviceRouter.regexpFilters.Delete(url)
			serviceRouter.removeFilterIndex(url)
		}
	}
}

// AddHandler 添加handler
// 全匹配和正则匹配分开存放, 正则表达式以':'符号开始, 如: /upload/:\S+
func (serviceRouter *ServiceRouter) AddHandler(method, url string, handler HandlerFunc) error {
	serviceRouter.checkInit()
	surl, err := buildHandlerURL(method, url)
	if nil != err {
		return err
	}
	if serviceRouter.isDebug {
		logs.Debugln("AddHandler:", surl)
	}
	if strings.Contains(url, ":") {
		serviceRouter.regexpHandlers.Put(surl, handler)
		serviceRouter.addHandlerIndex(surl)
	} else {
		serviceRouter.urlHandlers.Put(surl, handler)
	}
	return nil
}

// addHandlerIndex 添加handler索引, 安装路径长度排序
func (serviceRouter *ServiceRouter) addHandlerIndex(url string) {
	newArray := append(serviceRouter.regexpHandlersIndex, url)
	strutil.SortByLen(newArray, true)
	strutil.SortBySplitLen(newArray, "/", true)
	serviceRouter.regexpHandlersIndex = newArray
}

// removeHandlerIndex 删除handler索引, surl: 'POST /api'
func (serviceRouter *ServiceRouter) removeHandlerIndex(surl string) {
	if len(surl) > 0 {
		for i := 0; i < len(serviceRouter.regexpHandlersIndex); i++ {
			if serviceRouter.regexpHandlersIndex[i] == surl {
				serviceRouter.regexpHandlersIndex = append(serviceRouter.regexpHandlersIndex[:i], serviceRouter.regexpHandlersIndex[i+i:]...)
				break
			}
		}
	}
}

// RemoveHandler 删除一个路由表
func (serviceRouter *ServiceRouter) RemoveHandler(method, url string) {
	if len(url) == 0 {
		return
	}
	if serviceRouter.isDebug {
		logs.Debugln("RemoveHandler:", url)
	}
	surl, _ := buildHandlerURL(method, url)
	if nil != serviceRouter.regexpHandlers {
		if serviceRouter.regexpHandlers.ContainsKey(surl) {
			serviceRouter.regexpHandlers.Delete(surl)
			serviceRouter.removeHandlerIndex(surl)
		}
	}
	if nil != serviceRouter.urlHandlers {
		serviceRouter.urlHandlers.Delete(surl)
	}
}

// checkInit 检查是否初始化了
func (serviceRouter *ServiceRouter) checkInit() {
	if serviceRouter.initial {
		return
	} else {
		serviceRouter.initial = true
	}
	if nil == serviceRouter.regexpHandlers {
		serviceRouter.regexpHandlers = utypes.NewSafeMap()
		serviceRouter.regexpHandlersIndex = make([]string, 0)
	}
	if nil == serviceRouter.urlHandlers {
		serviceRouter.urlHandlers = utypes.NewSafeMap()
	}
	if nil == serviceRouter.regexpFilters {
		serviceRouter.regexpFilters = utypes.NewSafeMap()
	}
	if nil == serviceRouter.regexpFiltersIndex {
		serviceRouter.regexpFiltersIndex = make([]string, 0)
	}
	if nil == serviceRouter.runtimeError {
		serviceRouter.runtimeError = func(rw http.ResponseWriter, r *http.Request, err interface{}) {
			SendServerError(rw, fmt.Sprintf("%v", err))
		}
	}
}

// 格式method
func formatMethod(method string) string {
	method = strings.TrimSpace(method)
	if len(method) == 0 {
		return "ANY"
	}
	return strings.ToUpper(method)
}

// 拼接存储url, 格式: POST /api/:S+
func buildHandlerURL(method, url string) (string, error) {
	url = strings.TrimSpace(url)
	if len(url) == 0 {
		return "", errors.New("handler url is empty")
	}
	method = formatMethod(method)
	return method + " " + url, nil
}
