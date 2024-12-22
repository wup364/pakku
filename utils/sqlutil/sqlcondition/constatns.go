/// Copyright (C) 2024 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlcondition

import "github.com/wup364/pakku/utils/constants/sqlconditions"

// PaginationType 分页类型
type PaginationType int

const (
	// PageNumberType 页码模式
	PageNumberType PaginationType = iota
	// LimitOffsetType 限制偏移模式
	LimitOffsetType
)

// sqcSupportedGroupConditions 查询项与查询项之间支持的条件
var sqcSupportedGroupConditions = []string{sqlconditions.AND, sqlconditions.OR}

// sqcUnsupportedConditionErr 不支持的条件类型
var sqcUnsupportedConditionErr = "unsupported condition type '%s'"

// sqcSupportedItemConditions 查询项支持的条件表达式
var sqcSupportedItemConditions = map[string]struct{}{
	sqlconditions.EQUALS:             {},
	sqlconditions.LIKE:               {},
	sqlconditions.NOT:                {},
	sqlconditions.GREATER_THAN:       {},
	sqlconditions.GREATER_THAN_EQUAL: {},
	sqlconditions.LESS_THAN:          {},
	sqlconditions.LESS_THAN_EQUAL:    {},
	sqlconditions.IN:                 {},
	sqlconditions.BETWEEN:            {},
	sqlconditions.IS_NULL:            {},
	sqlconditions.IS_NOT_NULL:        {},
	// sqlconditions.EXISTS:      {},
	// sqlconditions.NOT_EXISTS":  {},
}
