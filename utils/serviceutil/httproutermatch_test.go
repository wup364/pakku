// Copyright (C) 2024 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package serviceutil

import (
	"testing"
)

// TestURLMatcher 测试主要的匹配功能
func TestURLMatcher(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
		params  map[string]string
	}{
		{
			name:    "精确匹配 - 成功",
			pattern: "/api/users",
			path:    "/api/users",
			want:    true,
			params:  map[string]string{},
		},
		{
			name:    "精确匹配 - 失败",
			pattern: "/api/users",
			path:    "/api/posts",
			want:    false,
			params:  map[string]string{},
		},
		{
			name:    "单层通配符 - 成功",
			pattern: "/api/:*",
			path:    "/api/users",
			want:    true,
			params:  map[string]string{},
		},
		{
			name:    "单层通配符 - 失败（多层）",
			pattern: "/api/:*",
			path:    "/api/users/123",
			want:    false,
			params:  map[string]string{},
		},
		{
			name:    "深度通配符 - 单层",
			pattern: "/api/:**",
			path:    "/api/users",
			want:    true,
			params:  map[string]string{},
		},
		{
			name:    "深度通配符 - 多层",
			pattern: "/api/:**",
			path:    "/api/users/123/posts",
			want:    true,
			params:  map[string]string{},
		},
		{
			name:    "多段匹配 - 成功",
			pattern: "/api/:*/posts/:*",
			path:    "/api/users/posts/123",
			want:    true,
			params:  map[string]string{},
		},
		{
			name:    "多段匹配 - 失败（段数不匹配）",
			pattern: "/api/:*/posts/:*",
			path:    "/api/users/posts/123/comments",
			want:    false,
			params:  map[string]string{},
		},
		{
			name:    "命名参数 - 单个",
			pattern: "/api/:id",
			path:    "/api/123",
			want:    true,
			params:  map[string]string{"id": "123"},
		},
		{
			name:    "命名参数 - 多个",
			pattern: "/api/:group/:id",
			path:    "/api/users/123",
			want:    true,
			params:  map[string]string{"group": "users", "id": "123"},
		},
		{
			name:    "混合模式",
			pattern: "/api/:group/posts/:id/comments",
			path:    "/api/users/posts/123/comments",
			want:    true,
			params:  map[string]string{"group": "users", "id": "123"},
		},
		{
			name:    "处理前后斜杠",
			pattern: "/api/users/",
			path:    "api/users",
			want:    true,
			params:  map[string]string{},
		},
		{
			name:    "空路径",
			pattern: "",
			path:    "/api/users",
			want:    false,
			params:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewURLMatcher(tt.pattern)

			// 测试匹配结果
			if got := matcher.Match(tt.path); got != tt.want {
				t.Errorf("URLMatcher.Match() = %v, want %v", got, tt.want)
			}

			// 测试参数提取
			if len(tt.params) > 0 {
				params := matcher.GetParams(tt.path)
				for k, v := range tt.params {
					if got, ok := params[k]; !ok || got != v {
						t.Errorf("URLMatcher.GetParams() key %s = %v, want %v", k, got, v)
					}
				}
			}
		})
	}
}

// TestURLMatcherEdgeCases 测试边缘情况
func TestURLMatcherEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{
			name:    "多余空格",
			pattern: "  /api/users  ",
			path:    "/api/users",
			want:    true,
		},
		{
			name:    "重复斜杠",
			pattern: "/api//users///",
			path:    "/api/users",
			want:    false,
		},
		{
			name:    "点号路径",
			pattern: "/api/:*",
			path:    "/api/.",
			want:    true,
		},
		{
			name:    "空段",
			pattern: "/api/:*/",
			path:    "/api//",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := NewURLMatcher(tt.pattern)
			if got := matcher.Match(tt.path); got != tt.want {
				t.Errorf("URLMatcher.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkURLMatcher 测试各种匹配模式的性能
func BenchmarkURLMatcher(b *testing.B) {
	benchmarks := []struct {
		name    string
		pattern string
		path    string
	}{
		{
			name:    "精确匹配",
			pattern: "/api/users",
			path:    "/api/users",
		},
		{
			name:    "单层通配符",
			pattern: "/api/:*",
			path:    "/api/users",
		},
		{
			name:    "深度通配符",
			pattern: "/api/:**",
			path:    "/api/users/123/posts/456/comments",
		},
		{
			name:    "多段匹配",
			pattern: "/api/:*/posts/:*/comments/:*",
			path:    "/api/users/posts/123/comments/456",
		},
		{
			name:    "命名参数",
			pattern: "/api/:group/:id",
			path:    "/api/users/123",
		},
		{
			name:    "混合模式",
			pattern: "/api/:group/posts/:id/comments",
			path:    "/api/users/posts/123/comments",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name+"_Match", func(b *testing.B) {
			matcher := NewURLMatcher(bm.pattern)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matcher.Match(bm.path)
			}
		})

		b.Run(bm.name+"_GetParams", func(b *testing.B) {
			matcher := NewURLMatcher(bm.pattern)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matcher.GetParams(bm.path)
			}
		})

		b.Run(bm.name+"_New", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				NewURLMatcher(bm.pattern)
			}
		})
	}
}

// BenchmarkURLMatcherParallel 测试并发性能
func BenchmarkURLMatcherParallel(b *testing.B) {
	patterns := []string{
		"/api/users",
		"/api/:*",
		"/api/:**",
		"/api/:*/posts/:*/comments/:*",
		"/api/:group/:id",
		"/api/:group/posts/:id/comments",
	}
	paths := []string{
		"/api/users",
		"/api/groups",
		"/api/users/123/posts/456/comments",
		"/api/users/posts/123/comments/456",
		"/api/users/123",
		"/api/users/posts/123/comments",
	}

	matchers := make([]*URLMatcher, len(patterns))
	for i, pattern := range patterns {
		matchers[i] = NewURLMatcher(pattern)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			idx := i % len(matchers)
			matchers[idx].Match(paths[idx])
			matchers[idx].GetParams(paths[idx])
			i++
		}
	})
}
