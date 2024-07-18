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

// debugLog 调试日志记录
func (srt *ServiceRouter) debugLog(msg ...interface{}) {
	if srt.isDebug {
		logs.Debugln(msg...)
	}
}

// ServiceRouter 根据注册的路由表调用对应的函数
// 优先匹配全url > 正则url > 默认处理器 > 404
func (srt *ServiceRouter) doHandle(w http.ResponseWriter, r *http.Request) {
	surl, err := buildHandlerURL(r.Method, r.URL.Path)
	if nil != err {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 1.0 如果是url全匹配, 则直接执行handler函数 - 有指定请求方式, POST, GET..
	if h, ok := srt.urlHandlers.Get(surl); ok {
		srt.debugLog("[URL.Handler.Path]", surl)
		h.(HandlerFunc)(w, r)
		return
	}

	// 1.1 如果是url全匹配, 则直接执行handler函数 - 无指定请求方式, POST, GET..
	anyurl, _ := buildHandlerURL("", r.URL.Path)
	if h, ok := srt.urlHandlers.Get(anyurl); ok {
		srt.debugLog("[URL.Handler.Path]", anyurl)
		h.(HandlerFunc)(w, r)
		return
	}

	// 2.0 如果是url正则检查, 则需要检查正则, 正则为':'后面的字符
	if srt.doExecExpHandler(surl, anyurl, w, r) {
		return
	}

	// 没有注册的地址, 使用默认处理器
	if srt.defaultHandler != nil {
		srt.debugLog("[URL.Handler.Default]", surl)
		srt.defaultHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// doExecExpHandle 执行正则、通配符规则, 有命中则返回true
func (srt *ServiceRouter) doExecExpHandler(surl string, anyurl string, w http.ResponseWriter, r *http.Request) bool {
	if len(srt.regexpHandlersIndex) == 0 {
		return false
	}

	var symbolIndex int
	for i := 0; i < len(srt.regexpHandlersIndex); i++ {
		if symbolIndex = strings.Index(srt.regexpHandlersIndex[i], ":"); symbolIndex == -1 {
			continue
		}
		if baseURL := srt.regexpHandlersIndex[i][:symbolIndex]; strings.HasPrefix(surl, baseURL) || strings.HasPrefix(anyurl, baseURL) {
			if matched, _ := regexp.MatchString(srt.getRegexpStr(srt.regexpHandlersIndex[i][symbolIndex+1:]), surl[symbolIndex:]); !matched {
				continue
			}
			if handler, ok := srt.regexpHandlers.Get(srt.regexpHandlersIndex[i]); ok {
				srt.debugLog("[URL.Handler.Regexp]", surl, srt.regexpHandlersIndex[i])
				handler.(HandlerFunc)(w, r)
				return true
			}
		}
	}
	return false
}

// getRegexpStr 获取正则字符串
func (srt *ServiceRouter) getRegexpStr(input string) string {
	if strutil.EqualsAny(input, "*", "*/") {
		return `^[^\s/]+$`

	} else if strutil.EqualsAny(input, "**", "**/") {
		return `^[^\s/][^\s]*$`

	} else {
		// 处理头部
		var tmp string
		if strings.HasPrefix(input, "**") {
			tmp = `^[^\s/][^\s]*` + input[2:]
		} else if strings.HasPrefix(input, "*") {
			tmp = `^[^\s/]+` + input[1:]
		}

		// 处理尾部
		if strings.HasSuffix(tmp, ":**") || strings.HasSuffix(tmp, ":**/") {
			tmp = tmp[:strings.LastIndex(tmp, ":**")] + `[^\s]*`
		} else if strings.HasSuffix(tmp, ":*") || strings.HasSuffix(tmp, ":*/") {
			tmp = tmp[:strings.LastIndex(tmp, ":*")] + `[^\s/]+`
		}

		// 处理中间部分
		tmp = strings.Replace(tmp, ":**", `[^\s]+`, -1)
		tmp = strings.Replace(tmp, ":*", `[^\s/]+`, -1)

		return tmp + "$"
	}
}

// doFilter 根据注册的过滤器表调用对应的函数
// 优先匹配全url > 正则url > 全局过滤器 > 直接通过
func (srt *ServiceRouter) doFilter(w http.ResponseWriter, r *http.Request) {
	// 1. 执行正则url过滤器, 有命中则返回true
	if srt.doExecuteExpFilter(w, r) {
		return
	}

	// 2. 检查是否有全局过滤器存在, 如果有则执行它
	if nil != srt.defaultFileter {
		srt.debugLog("[URL.Filter.Default]", r.URL.Path)
		if srt.defaultFileter(w, r) {
			srt.doHandle(w, r)
		}
		return
	}

	// 3. 啥也没有设定
	srt.doHandle(w, r)
}

// doExecuteExpFilter 执行正则url过滤器, 有命中则返回true
func (srt *ServiceRouter) doExecuteExpFilter(w http.ResponseWriter, r *http.Request) bool {
	if matched := srt.getMatchedFilter(r.URL.Path); len(matched) > 0 {
		for i := 0; i < len(matched); i++ {
			if h, ok := srt.regexpFilters.Get(matched[i]); ok {
				srt.debugLog("[URL.Filter]", matched[i])
				if !h.(FilterFunc)(w, r) {
					return true
				}
			}
		}

		// 符合所有过滤器要求
		srt.doHandle(w, r)
		return true
	}
	return false
}

// getMatchedFilter 获取匹配的过滤器
func (srt *ServiceRouter) getMatchedFilter(urlPath string) []string {
	if len(srt.regexpFiltersIndex) == 0 || len(urlPath) == 0 {
		return nil
	}

	matched := make([]string, 0)
	for i := 0; i < len(srt.regexpFiltersIndex); i++ {
		var symbolIndex int
		if symbolIndex = strings.Index(srt.regexpFiltersIndex[i], ":"); symbolIndex == -1 {
			if urlPath == srt.regexpFiltersIndex[i] {
				matched = append(matched, srt.regexpFiltersIndex[i])
			}
			continue
		}

		// 正则匹配
		if baseURL := srt.regexpFiltersIndex[i][:symbolIndex]; strings.HasPrefix(urlPath, baseURL) {
			if ok, _ := regexp.MatchString(srt.getRegexpStr(srt.regexpFiltersIndex[i][symbolIndex+1:]), urlPath[symbolIndex:]); ok {
				matched = append(matched, srt.regexpFiltersIndex[i])
			}
		}
	}
	return matched
}

// ServeHTTP 实现http.Server接口的ServeHTTP方法
func (srt *ServiceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if srt.runtimeError != nil {
			if err := recover(); nil != err {
				srt.runtimeError(w, r, err)
			}
		}
	}()

	srt.checkInit()
	srt.doFilter(w, r)
}

// ClearHandlersMap 清空路由表
func (srt *ServiceRouter) ClearHandlersMap() {
	srt.urlHandlers.Clear()
	srt.regexpHandlers.Clear()
	srt.regexpHandlersIndex = make([]string, 0)
}

// SetDebug 是否输出url请求信息
func (srt *ServiceRouter) SetDebug(isDebug bool) {
	srt.isDebug = isDebug
}

// SetDefaultHandler 设置默认响应函数, 当无匹配时触发
func (srt *ServiceRouter) SetDefaultHandler(defaultHandler HandlerFunc) {
	logs.Debugln("The default handler has been set")
	srt.defaultHandler = defaultHandler
}

// SetRuntimeErrorHandler 设置全局handler执行异常捕获
func (srt *ServiceRouter) SetRuntimeErrorHandler(h RuntimeErrorHandlerFunc) {
	srt.runtimeError = h
}

// SetDefaultFilter 设置默认过滤器, 设置后, 如果不调用next函数则不进行下一步处理
// type FilterFunc func(http.ResponseWriter, *http.Request, func( ))
func (srt *ServiceRouter) SetDefaultFilter(globalFilter FilterFunc) {
	logs.Debugln("The default filter has been set")
	srt.defaultFileter = globalFilter
}

// AddURLFilter 设置url过滤器, 设置后, 如果不调用next函数则不进行下一步处理
// 过滤器有优先调用权, 正则匹配路径有先后顺序
// type FilterFunc func(http.ResponseWriter, *http.Request, func( ))
func (srt *ServiceRouter) AddURLFilter(url string, filter FilterFunc) error {
	if len(url) == 0 {
		return errors.New("filter url is empty")
	}
	srt.checkInit()
	url = parse2UnixPath(url)
	logs.Debugln("AddURLFilter: ", url)
	srt.regexpFilters.Put(url, filter)
	srt.addFilterIndex(url)
	return nil
}

// addFilterIndex 添加filter索引, 安装路径长度排序
func (srt *ServiceRouter) addFilterIndex(url string) {
	newArray := append(srt.regexpFiltersIndex, url)
	strutil.SortByLen(newArray, true)
	strutil.SortBySplitLen(newArray, "/", true)
	srt.regexpFiltersIndex = newArray
}

// removeFilterIndex 删除filter索引
func (srt *ServiceRouter) removeFilterIndex(url string) {
	if len(url) > 0 {
		for i := 0; i < len(srt.regexpFiltersIndex); i++ {
			if srt.regexpFiltersIndex[i] == url {
				srt.regexpFiltersIndex = append(srt.regexpFiltersIndex[:i], srt.regexpFiltersIndex[i+i:]...)
				break
			}
		}
	}
}

// RemoveFilter 删除一个过滤器
func (srt *ServiceRouter) RemoveFilter(url string) {
	if len(url) == 0 {
		return
	}
	logs.Debugln("RemoveFilter:", url)
	if nil != srt.regexpFilters && srt.regexpFilters.ContainsKey(url) {
		srt.regexpFilters.Delete(url)
		srt.removeFilterIndex(url)
	}
}

// AddHandler 添加handler
// 全匹配和正则匹配分开存放, 正则表达式以':'符号开始, 如: /upload/:\S+
func (srt *ServiceRouter) AddHandler(method, url string, handler HandlerFunc) error {
	srt.checkInit()
	if surl, err := buildHandlerURL(method, parse2UnixPath(url)); nil != err {
		return err
	} else {
		logs.Debugln("AddHandler:", surl)
		if strings.Contains(url, ":") {
			srt.regexpHandlers.Put(surl, handler)
			srt.addHandlerIndex(surl)
		} else {
			srt.urlHandlers.Put(surl, handler)
		}
	}
	return nil
}

// addHandlerIndex 添加handler索引, 安装路径长度排序
func (srt *ServiceRouter) addHandlerIndex(url string) {
	newArray := append(srt.regexpHandlersIndex, url)
	strutil.SortByLen(newArray, true)
	strutil.SortBySplitLen(newArray, "/", true)
	srt.regexpHandlersIndex = newArray
}

// removeHandlerIndex 删除handler索引, surl: 'POST /api'
func (srt *ServiceRouter) removeHandlerIndex(surl string) {
	if len(surl) > 0 {
		for i := 0; i < len(srt.regexpHandlersIndex); i++ {
			if srt.regexpHandlersIndex[i] == surl {
				srt.regexpHandlersIndex = append(srt.regexpHandlersIndex[:i], srt.regexpHandlersIndex[i+i:]...)
				break
			}
		}
	}
}

// RemoveHandler 删除一个路由表
func (srt *ServiceRouter) RemoveHandler(method, url string) {
	if len(url) == 0 {
		return
	}
	logs.Debugln("RemoveHandler:", url)
	surl, _ := buildHandlerURL(method, url)
	if nil != srt.regexpHandlers && srt.regexpHandlers.ContainsKey(surl) {
		srt.regexpHandlers.Delete(surl)
		srt.removeHandlerIndex(surl)
	}
	if nil != srt.urlHandlers {
		srt.urlHandlers.Delete(surl)
	}
}

// checkInit 检查是否初始化了
func (srt *ServiceRouter) checkInit() {
	if srt.initial {
		return
	}
	srt.initial = true

	if nil == srt.regexpHandlers {
		srt.regexpHandlers = utypes.NewSafeMap()
		srt.regexpHandlersIndex = make([]string, 0)
	}
	if nil == srt.urlHandlers {
		srt.urlHandlers = utypes.NewSafeMap()
	}
	if nil == srt.regexpFilters {
		srt.regexpFilters = utypes.NewSafeMap()
	}
	if nil == srt.regexpFiltersIndex {
		srt.regexpFiltersIndex = make([]string, 0)
	}
	if nil == srt.runtimeError {
		srt.runtimeError = func(rw http.ResponseWriter, r *http.Request, err interface{}) {
			SendServerError(rw, fmt.Sprintf("%v", err))
			logs.Errorln(err)
		}
	}
}

// Parse2UnixPath 格式化入参路径
func parse2UnixPath(url string) string {
	if i := strings.Index(url, ":"); i > 0 {
		return strutil.Parse2UnixPath(url[:i+1]) + url[i+1:]
	} else {
		url = strutil.Parse2UnixPath(url)
	}
	return url
}

// 格式method
func formatMethod(method string) string {
	if method = strings.TrimSpace(method); len(method) == 0 {
		return "ANY"
	}
	return strings.ToUpper(method)
}

// 拼接存储url, 格式: POST /api/:S+
func buildHandlerURL(method, url string) (string, error) {
	if url = strings.TrimSpace(url); len(url) == 0 {
		return "", errors.New("handler url is empty")
	}
	method = formatMethod(method)
	return method + " " + url, nil
}
