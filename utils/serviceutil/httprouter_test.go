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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceRouter(t *testing.T) {
	// 创建新的路由器实例
	router := NewServiceRouter()
	router.SetDebug(true)

	// 测试基本路由处理
	t.Run("Basic Route Handling", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("success"))
		}

		err := router.AddHandler("GET", "/test", handler)
		if err != nil {
			t.Errorf("添加处理器失败: %v", err)
		}

		// 创建测试请求
		req, _ := http.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("处理器返回了错误的状态码: got %v want %v", rr.Code, http.StatusOK)
		}

		if rr.Body.String() != "success" {
			t.Errorf("处理器返回了错误的响应体: got %v want %v", rr.Body.String(), "success")
		}
	})

	// 测试过滤器
	t.Run("Filter Chain", func(t *testing.T) {
		filterCalled := false
		filter := func(w http.ResponseWriter, r *http.Request) bool {
			filterCalled = true
			return true
		}

		err := router.AddURLFilter("/filtered", filter)
		if err != nil {
			t.Errorf("添加过滤器失败: %v", err)
		}

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("filtered"))
		}

		router.AddHandler("GET", "/filtered", handler)

		req, _ := http.NewRequest("GET", "/filtered", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if !filterCalled {
			t.Error("过滤器未被调用")
		}

		if rr.Body.String() != "filtered" {
			t.Errorf("处理器返回了错误的响应体: got %v want %v", rr.Body.String(), "filtered")
		}
	})

	// 测试通配符路由
	t.Run("Wildcard Routes", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("wildcard"))
		}

		router.AddHandler("GET", "/test/:*", handler)

		req, _ := http.NewRequest("GET", "/test/anything", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Body.String() != "wildcard" {
			t.Errorf("通配符路由处理失败: got %v want %v", rr.Body.String(), "wildcard")
		}
	})

	// 测试默认处理器
	t.Run("Default Handler", func(t *testing.T) {
		defaultHandler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}

		router.SetDefaultHandler(defaultHandler)

		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("默认处理器返回了错误的状态码: got %v want %v", rr.Code, http.StatusNotFound)
		}
	})

	// 测试移除路由
	t.Run("Remove Handler", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("temp"))
		}

		router.AddHandler("GET", "/temp", handler)
		router.RemoveHandler("GET", "/temp")

		req, _ := http.NewRequest("GET", "/temp", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("移除处理器后仍然可以访问: got %v want %v", rr.Code, http.StatusNotFound)
		}
	})
}

// 测试运行时错误处理
func TestRuntimeErrorHandler(t *testing.T) {
	router := NewServiceRouter()

	errorHandlerCalled := false
	router.SetRuntimeErrorHandler(func(w http.ResponseWriter, r *http.Request, err any) {
		errorHandlerCalled = true
		w.WriteHeader(http.StatusInternalServerError)
	})

	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("测试错误")
	}

	router.AddHandler("GET", "/error", handler)

	req, _ := http.NewRequest("GET", "/error", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if !errorHandlerCalled {
		t.Error("运行时错误处理器未被调用")
	}

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("错误处理器返回了错误的状态码: got %v want %v", rr.Code, http.StatusInternalServerError)
	}
}

// 性能测试
func BenchmarkRouter(b *testing.B) {
	router := NewServiceRouter()

	// 注册一些测试路由
	router.AddHandler("GET", "/static", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	router.AddHandler("GET", "/user/:id", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	router.AddHandler("GET", "/api/:version/users/:id/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	router.AddHandler("GET", "/static/:*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// 测试静态路由性能
	b.Run("Static Route", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "/static", nil)
		rr := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			router.ServeHTTP(rr, req)
		}
	})

	// 测试参数路由性能
	b.Run("Parameterized Route", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "/user/123", nil)
		rr := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			router.ServeHTTP(rr, req)
		}
	})

	// 测试复杂嵌套参数路由性能
	b.Run("Complex Parameterized Route", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "/api/v1/users/123/profile", nil)
		rr := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			router.ServeHTTP(rr, req)
		}
	})

	// 测试通配符路由性能
	b.Run("Wildcard Route", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "/static/css/style.css", nil)
		rr := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			router.ServeHTTP(rr, req)
		}
	})

	// 测试不存在的路由性能
	b.Run("Non-Existent Route", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "/not/found", nil)
		rr := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			router.ServeHTTP(rr, req)
		}
	})
}
