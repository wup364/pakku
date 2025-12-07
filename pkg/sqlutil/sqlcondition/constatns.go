// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlcondition

import "github.com/wup364/pakku/pkg/constants/sqlconditions"

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
