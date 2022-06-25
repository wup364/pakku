// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package serviceutil

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
)

func TestStartService(t *testing.T) {
	err := startService("127.0.0.1:8080", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		// w.Write([]byte("default" + r.URL.String()))
	})
	if nil != err {
		t.Error(err)
	}
}

// StartService 注册路由,并启动服务
// 此函数为ServiceRouter简版的使用方式, 可以根据ServiceRouter自己实现
// 默认读写都不设置超时
// isDebug 参数在生产环境注意关闭, 有一倍性能差距 15018/sec-> 30173/sec
func startService(addr string, defaultHandler HandlerFunc) error {
	if addr == "" {
		return errors.New("addr is nil or empty")
	}
	router := &ServiceRouter{}
	router.isDebug = true
	router.defaultHandler = defaultHandler
	AddURLHandlers(router)
	AddURLFilters(router)

	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    0,
		WriteTimeout:   0,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}

func AddURLHandlers(router *ServiceRouter) {
	// 100 个无效测试数据, 模拟url查找
	for i := 0; i < 100; i++ {
		router.AddHandler(http.MethodPost, "/hello"+strconv.Itoa(i)+"/:[0-9]", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello/:[0-9] : " + r.Method + ": " + r.URL.String()))
		})
		router.AddHandler(http.MethodGet, "/hello"+strconv.Itoa(i), func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello : " + r.Method + ": " + r.URL.String()))
		})
	}
	router.AddHandler(http.MethodGet, "/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello : " + r.Method + ": " + r.URL.String()))
	})
	router.AddHandler(http.MethodPost, "/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello : " + r.Method + ": " + r.URL.String()))
	})
	router.AddHandler(http.MethodGet, "/hello/:[0-9]", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello/:[0-9] : " + r.Method + ": " + r.URL.String()))
	})
	router.AddHandler(http.MethodPost, "/hello/:[0-9]", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello/:[0-9] : " + r.Method + ": " + r.URL.String()))
	})
	router.AddHandler("", "/hello/:[a-zA-Z]", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello/:[a-zA-Z] : " + r.Method + ": " + r.URL.String()))
	})
}

func AddURLFilters(router *ServiceRouter) {
	router.AddURLFilter("/", func(rw http.ResponseWriter, r *http.Request) bool {
		return true
	})
	router.AddURLFilter("/hello/xxxxx", func(rw http.ResponseWriter, r *http.Request) bool {
		return true
	})
	router.AddURLFilter("/hello/xx/x", func(rw http.ResponseWriter, r *http.Request) bool {
		SendServerError(rw, "forbid")
		return false
	})
	router.AddURLFilter("/hello/xxx", func(rw http.ResponseWriter, r *http.Request) bool {
		return true
	})
	router.AddURLFilter("/hello:*", func(rw http.ResponseWriter, r *http.Request) bool {
		return true
	})
	router.AddURLFilter("/hello/:[a-zA-Z]", func(rw http.ResponseWriter, r *http.Request) bool {
		return true
	})
	router.AddURLFilter("/hello", func(rw http.ResponseWriter, r *http.Request) bool {
		return true
	})
}
