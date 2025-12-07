// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package serviceutil

// http服务器工具-URL路由管理
// ServiceRouter 实现了 http.Server 接口的 ServeHTTP 方法, 提供 URL 路由管理功能。
// 请求处理逻辑: 过滤器 > 路径匹配 > END

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/wup364/pakku/pkg/logs"
	"github.com/wup364/pakku/pkg/strutil"
	"github.com/wup364/pakku/pkg/utypes"
)

// contextKey 自定义类型
type contextKey string

// urlParamsKey 上下文key
const urlParamsKey contextKey = "urlParams"

// HandlerFunc 定义请求处理器
type HandlerFunc func(http.ResponseWriter, *http.Request)

// FilterFunc http请求过滤器, 返回bool, true: 继续, false: 停止
type FilterFunc func(http.ResponseWriter, *http.Request) bool

// RuntimeErrorHandlerFunc 未知异常处理函数
type RuntimeErrorHandlerFunc func(http.ResponseWriter, *http.Request, any)

// routerKey 自定义类型
type routerKey string

// NewServiceRouter New service router
func NewServiceRouter() (router *ServiceRouter) {
	router = &ServiceRouter{}
	router.initializeServiceRouter()
	return router
}

// ServiceRouter 实现了 http.Server 接口的 ServeHTTP 方法, 提供 URL 路由管理功能。
// 请求处理逻辑: 过滤器 > 路径匹配 > END
type ServiceRouter struct {
	initial         bool                                        // 是否已初始化
	isDebug         bool                                        // 调试模式, 启用时输出详细日志
	enableURLParam  bool                                        // 是否启用URL参数注入到上下文
	urlHandlers     *utypes.SafeMap[routerKey, *HandlerEntry]   // URL处理器映射表
	urlFilters      *utypes.SafeMap[routerKey, *URLFilterEntry] // URL过滤器映射表
	handlersIndex   []string                                    // 处理器索引, 保持注册顺序
	urlFiltersIndex []string                                    // 过滤器索引, 保持注册顺序
	defaultFileter  FilterFunc                                  // 默认过滤器
	defaultHandler  HandlerFunc                                 // 默认处理器, 处理未匹配的请求
	runtimeError    RuntimeErrorHandlerFunc                     // 运行时错误处理函数
}

// URLFilterEntry url过滤器结构体
type URLFilterEntry struct {
	matcher *URLMatcher
	filter  FilterFunc
}

// HandlerEntry 添加新的结构体定义
type HandlerEntry struct {
	matcher *URLMatcher
	handler HandlerFunc
	method  string
}

// ServeHTTP 实现 http.Handler 接口, 处理所有 HTTP 请求
func (srt *ServiceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); nil != err {
			if srt.runtimeError != nil {
				srt.runtimeError(w, r, err)
			} else {
				logs.Error(err)
			}
		}
	}()

	srt.doFilter(w, r)
}

// SetDebug 设置是否启用调试模式, 启用后会输出详细的请求处理日志
func (srt *ServiceRouter) SetDebug(isDebug bool) {
	srt.isDebug = isDebug
}

// SetDefaultHandler 设置默认的请求处理器, 当没有匹配的路由时被调用
func (srt *ServiceRouter) SetDefaultHandler(defaultHandler HandlerFunc) {
	logs.Debug("The default handler has been set")
	srt.defaultHandler = defaultHandler
}

// SetRuntimeErrorHandler 设置运行时错误处理函数, 用于处理请求处理过程中的panic
func (srt *ServiceRouter) SetRuntimeErrorHandler(h RuntimeErrorHandlerFunc) {
	srt.runtimeError = h
}

// SetDefaultFilter 设置默认的请求过滤器, 在所有特定路由过滤器之后执行
func (srt *ServiceRouter) SetDefaultFilter(globalFilter FilterFunc) {
	logs.Debug("The default filter has been set")
	srt.defaultFileter = globalFilter
}

// AddURLFilter 添加URL过滤器, 如: /api/**
func (srt *ServiceRouter) AddURLFilter(url string, filter FilterFunc) error {
	url = strutil.Parse2UnixPath(url)
	if len(url) == 0 {
		return errors.New("filter url is empty")
	} else {
		logs.Debug("AddURLFilter: ", url)
	}

	entry := &URLFilterEntry{
		matcher: NewURLMatcher(url),
		filter:  filter,
	}
	srt.urlFilters.Put(routerKey(url), entry)
	srt.appendAndSortFilterIndex(url)
	return nil
}

// AddHandler 添加URL处理器, 如: POST /api/:*, GET /api/:id
func (srt *ServiceRouter) AddHandler(method, url string, handler HandlerFunc) error {
	surl, err := srt.buildHandlerURL(method, url)
	if err != nil {
		return err
	}
	logs.Debug("AddHandler:", surl)

	entry := &HandlerEntry{
		matcher: NewURLMatcher(url),
		handler: handler,
		method:  srt.formatMethod(method),
	}

	srt.urlHandlers.Put(routerKey(surl), entry)
	srt.appendAndSortHandlerIndex(surl)
	return nil
}

// RemoveFilter 移除指定URL的过滤器
func (srt *ServiceRouter) RemoveFilter(url string) {
	if len(url) == 0 {
		return
	}
	logs.Debug("RemoveFilter:", url)
	if nil != srt.urlFilters && srt.urlFilters.ContainsKey(routerKey(url)) {
		srt.urlFilters.Delete(routerKey(url))
		srt.deleteFilterIndex(url)
	}
}

// RemoveHandler 移除指定URL和方法的处理器
func (srt *ServiceRouter) RemoveHandler(method, url string) {
	surl, err := srt.buildHandlerURL(method, url)
	if err != nil {
		return
	}
	logs.Debug("RemoveHandler:", surl)

	if srt.urlHandlers.ContainsKey(routerKey(surl)) {
		srt.urlHandlers.Delete(routerKey(surl))
		srt.deleteHandlerIndex(surl)
	}
}

// ClearHandlersMap 清空所有注册的处理器
func (srt *ServiceRouter) ClearHandlersMap() {
	srt.urlHandlers.Clear()
	srt.handlersIndex = make([]string, 0)
}

// doFilter 使用URLMatcher重新实现
func (srt *ServiceRouter) doFilter(w http.ResponseWriter, r *http.Request) {
	// 1. 执行URL过滤器匹配
	if srt.doExecuteURLFilter(w, r) {
		return
	}

	// 2. 检查是否有全局过滤器存在
	if nil != srt.defaultFileter {
		srt.debugLog("[URL.Filter.Default]", r.URL.Path)
		if srt.defaultFileter(w, r) {
			srt.doHandle(w, r)
		}
		return
	}

	// 3. 无过滤器情况
	srt.doHandle(w, r)
}

// doExecuteURLFilter 用URLMatcher执行过滤器匹配
func (srt *ServiceRouter) doExecuteURLFilter(w http.ResponseWriter, r *http.Request) bool {
	path := r.URL.Path

	// 按照注册顺序遍历所有过滤器
	for _, pattern := range srt.urlFiltersIndex {
		if entry, ok := srt.urlFilters.Get(routerKey(pattern)); ok {
			if entry.matcher.Match(path) {
				srt.debugLog("[URL.Filter]", pattern)
				if !entry.filter(w, r) {
					return true
				}
			}
		}
	}

	// 所有匹配的过滤器都通过
	srt.doHandle(w, r)
	return true
}

// ServiceRouter 根据注册的路由表调用对应的函数, 优先级: 匹配url > 默认处理器 > 404
func (srt *ServiceRouter) doHandle(w http.ResponseWriter, r *http.Request) {
	if entry := srt.findMatchingHandler(r); entry != nil {
		srt.executeHandler(w, r, entry)
		return
	}

	// 使用默认处理器
	if srt.defaultHandler != nil {
		srt.debugLog("[URL.Handler.Default]", r.URL.Path)
		srt.defaultHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// findMatchingHandler 查找匹配的处理器
func (srt *ServiceRouter) findMatchingHandler(r *http.Request) *HandlerEntry {
	for _, pattern := range srt.handlersIndex {
		if entry, ok := srt.urlHandlers.Get(routerKey(pattern)); ok {
			if entry.method != "ANY" && entry.method != r.Method {
				continue
			}
			if entry.matcher.Match(r.URL.Path) {
				srt.debugLog("[URL.Handler]", pattern)
				return entry
			}
		}
	}
	return nil
}

// executeHandler 执行处理器
func (srt *ServiceRouter) executeHandler(w http.ResponseWriter, r *http.Request, entry *HandlerEntry) {
	if srt.enableURLParam {
		params := entry.matcher.GetParams(r.URL.Path)
		ctx := context.WithValue(r.Context(), urlParamsKey, params)
		entry.handler(w, r.WithContext(ctx))
	} else {
		entry.handler(w, r)
	}
}

// appendAndSortHandlerIndex 添加handler索引, 按照路径层级深度排序
func (srt *ServiceRouter) appendAndSortHandlerIndex(url string) {
	// 添加新的URL到索引中
	newArray := append(srt.handlersIndex, url)
	srt.sortPathArray(newArray, false)
	srt.handlersIndex = newArray
}

// appendAndSortFilterIndex 添加filter索引, 按照路径层级深度排序
func (srt *ServiceRouter) appendAndSortFilterIndex(url string) {
	newArray := append(srt.urlFiltersIndex, url)
	srt.sortPathArray(newArray, true)
	srt.urlFiltersIndex = newArray
}

// sortPathArray 按照路径层级深度排序
func (srt *ServiceRouter) sortPathArray(newArray []string, isFilter bool) {
	sort.Slice(newArray, func(i, j int) bool {
		pathI, pathJ := newArray[i], newArray[j]

		if !isFilter {
			// 提取路径部分（去掉HTTP方法）
			pathI = strings.SplitN(newArray[i], " ", 2)[1]
			pathJ = strings.SplitN(newArray[j], " ", 2)[1]
		}

		// 计算路径层级深度
		pathsI := strings.Split(strings.Trim(pathI, "/"), "/")
		pathsJ := strings.Split(strings.Trim(pathJ, "/"), "/")

		// 根据层级排序
		if depthI, depthJ := len(pathsI), len(pathsJ); depthI == depthJ {
			patternI := strings.Count(pathI, ":")
			patternJ := strings.Count(pathJ, ":")
			if patternI == patternJ {
				tmpPathI := strings.ReplaceAll(strings.ReplaceAll(pathI, ":", ""), "*", "")
				tmpPathJ := strings.ReplaceAll(strings.ReplaceAll(pathJ, ":", ""), "*", "")
				if lenI, lenJ := len(tmpPathI), len(tmpPathJ); lenI == lenJ {
					return strings.Count(pathI, "*") < strings.Count(pathJ, "*")
				} else {
					return lenI > lenJ
				}
			}

			return patternI < patternJ

			// 按照层级深度从小到大排序
		} else {
			return depthI < depthJ
		}
	})
}

// deleteFilterIndex 删除filter索引
func (srt *ServiceRouter) deleteFilterIndex(url string) {
	if len(url) > 0 {
		for i := 0; i < len(srt.urlFiltersIndex); i++ {
			if srt.urlFiltersIndex[i] == url {
				srt.urlFiltersIndex = append(srt.urlFiltersIndex[:i], srt.urlFiltersIndex[i+i:]...)
				break
			}
		}
	}
}

// deleteHandlerIndex 删除handler索引
func (srt *ServiceRouter) deleteHandlerIndex(url string) {
	if len(url) > 0 {
		for i := 0; i < len(srt.handlersIndex); i++ {
			if srt.handlersIndex[i] == url {
				srt.handlersIndex = append(srt.handlersIndex[:i], srt.handlersIndex[i+1:]...)
				break
			}
		}
	}
}

// initializeServiceRouter 修改初始化方法
func (srt *ServiceRouter) initializeServiceRouter() {
	if srt.initial {
		return
	}
	srt.initial = true
	srt.enableURLParam = true

	if nil == srt.urlHandlers {
		srt.urlHandlers = utypes.NewSafeMap[routerKey, *HandlerEntry]()
		srt.handlersIndex = make([]string, 0)
	}
	if nil == srt.urlFilters {
		srt.urlFilters = utypes.NewSafeMap[routerKey, *URLFilterEntry]()
	}
	if nil == srt.urlFiltersIndex {
		srt.urlFiltersIndex = make([]string, 0)
	}
	if nil == srt.runtimeError {
		srt.runtimeError = func(rw http.ResponseWriter, r *http.Request, err any) {
			SendServerError(rw, fmt.Sprintf("%v", err))
			logs.Error(err)
		}
	}
}

// debugLog 调试日志记录
func (srt *ServiceRouter) debugLog(msg ...any) {
	if srt.isDebug {
		logs.Debug(msg...)
	}
}

// formatMethod 格式method
func (srt *ServiceRouter) formatMethod(method string) string {
	if method = strings.TrimSpace(method); len(method) == 0 {
		return "ANY"
	}
	return strings.ToUpper(method)
}

// buildHandlerURL 拼接存储url, 格式: POST /api/:*
func (srt *ServiceRouter) buildHandlerURL(method, url string) (string, error) {
	if url = strings.TrimSpace(url); len(url) == 0 {
		return "", errors.New("handler url is empty")
	}
	method = srt.formatMethod(method)
	return method + " " + url, nil
}

// EnableURLParam 设置是否启用URL参数注入到请求上下文中
func (srt *ServiceRouter) EnableURLParam(enable bool) {
	srt.enableURLParam = enable
}

// GetURLParam 从请求上下文中获取URL参数值
// paramName: 参数名
// 返回参数值, 如果参数不存在则返回空字符串
func GetURLParam(r *http.Request, paramName string) string {
	if params, ok := r.Context().Value(urlParamsKey).(map[string]string); ok {
		return params[paramName]
	}
	return ""
}
