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
	"strings"
)

// URLMatcher URL路由匹配器
type URLMatcher struct {
	pattern string   // 标准化后的匹配模式
	parts   []string // 预分割的模式部分
}

// NewURLMatcher 创建新的URL匹配器
func NewURLMatcher(pattern string) *URLMatcher {
	pattern = strings.Trim(strings.TrimSpace(pattern), "/")
	return &URLMatcher{
		pattern: pattern,
		parts:   strings.Split(pattern, "/"),
	}
}

// Match 检查路径是否匹配当前模式
func (m *URLMatcher) Match(path string) bool {
	// 处理空路径
	if len(path) == 0 || len(m.pattern) == 0 {
		return false
	}

	// 标准化路径
	path = strings.Trim(path, "/")

	// 处理精确匹配
	if !strings.Contains(m.pattern, ":") {
		return m.pattern == path
	}

	// 分割路径
	pathParts := strings.Split(path, "/")

	// 处理深度通配符 :**
	if m.hasDeepWildcard(m.parts) {
		return m.matchDeepWildcard(pathParts, m.parts)
	}

	// 处理普通匹配
	return m.matchParts(pathParts, m.parts)
}

// matchDeepWildcard 处理包含 :** 的匹配
func (m *URLMatcher) matchDeepWildcard(pathParts, patternParts []string) bool {
	// 找到 :** 的位置
	deepWildcardIndex := -1
	for i, part := range patternParts {
		if part == ":**" {
			deepWildcardIndex = i
			break
		}
	}

	// 检查 :** 之前的部分
	if deepWildcardIndex == -1 || len(pathParts) < deepWildcardIndex {
		return false
	}

	// 匹配 :** 之前的部分
	for i := 0; i < deepWildcardIndex; i++ {
		if !m.matchPart(pathParts[i], patternParts[i]) {
			return false
		}
	}

	return true
}

// matchParts 匹配路径段
func (m *URLMatcher) matchParts(pathParts, patternParts []string) bool {
	// 如果段数不同，则不匹配
	if len(pathParts) != len(patternParts) {
		return false
	}

	// 逐段匹配
	for i := 0; i < len(patternParts); i++ {
		if !m.matchPart(pathParts[i], patternParts[i]) {
			return false
		}
	}

	return true
}

// matchPart 匹配单个路径段
func (m *URLMatcher) matchPart(pathPart, patternPart string) bool {
	// 处理通配符
	if strings.HasPrefix(patternPart, ":") {
		return true // :* 或 :id 都匹配任意单个段
	}

	// 精确匹配
	return pathPart == patternPart
}

// hasDeepWildcard 检查是否包含深度通配符
func (m *URLMatcher) hasDeepWildcard(parts []string) bool {
	for _, part := range parts {
		if part == ":**" {
			return true
		}
	}
	return false
}

// GetParams 获取URL中的参数
func (m *URLMatcher) GetParams(path string) map[string]string {
	params := make(map[string]string)

	// 标准化路径
	path = strings.Trim(path, "/")
	pathParts := strings.Split(path, "/")

	// 提取参数
	for i := 0; i < len(m.parts) && i < len(pathParts); i++ {
		if strings.HasPrefix(m.parts[i], ":") {
			paramName := strings.TrimPrefix(m.parts[i], ":")
			if paramName == "*" || paramName == "**" {
				continue // 跳过通配符
			}
			params[paramName] = pathParts[i]
		}
	}

	return params
}
