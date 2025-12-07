// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

// 简单查询条件组合sql生成器
package sqlcondition

import (
	"fmt"
	"strings"

	"github.com/wup364/pakku/pkg/constants/sqlconditions"
)

// NewQueryConditionBuilder 创建一个新的查询构建器
func NewQueryConditionBuilder() *QueryConditionBuilder {
	return &QueryConditionBuilder{
		condition: sqlconditions.AND,
	}
}

// QueryConditionBuilder 查询条件构建器
type QueryConditionBuilder struct {
	condition    string
	groups       []QueryConditionGroup
	currentGroup *QueryConditionGroup
	orderBy      []OrderBy
	pagination   *Pagination
	dbType       string
}

// SetDBType 设置数据库类型
func (b *QueryConditionBuilder) SetDBType(dbType string) *QueryConditionBuilder {
	b.dbType = dbType
	return b
}

// OrderByAsc 添加升序排序
func (b *QueryConditionBuilder) OrderByAsc(fields ...string) *QueryConditionBuilder {
	for _, field := range fields {
		b.orderBy = append(b.orderBy, OrderBy{Field: field, Order: "ASC"})
	}
	return b
}

// OrderByDesc 添加降序排序
func (b *QueryConditionBuilder) OrderByDesc(fields ...string) *QueryConditionBuilder {
	for _, field := range fields {
		b.orderBy = append(b.orderBy, OrderBy{Field: field, Order: "DESC"})
	}
	return b
}

// OrderBy 添加自定义排序
func (b *QueryConditionBuilder) OrderBy(orders map[string]string) *QueryConditionBuilder {
	for field, order := range orders {
		upperOrder := strings.ToUpper(order)
		if upperOrder != "ASC" && upperOrder != "DESC" {
			upperOrder = "ASC" // 默认使用升序
		}
		b.orderBy = append(b.orderBy, OrderBy{Field: field, Order: upperOrder})
	}
	return b
}

// SetPageNumberPagination 设置页码模式分页
func (b *QueryConditionBuilder) SetPageNumberPagination(pageSize, pageNumber int) *QueryConditionBuilder {
	b.pagination = &Pagination{
		Type:       PageNumberType,
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}
	return b
}

// SetLimitOffsetPagination 设置限制偏移模式分页
func (b *QueryConditionBuilder) SetLimitOffsetPagination(limit, offset int) *QueryConditionBuilder {
	b.pagination = &Pagination{
		Type:     LimitOffsetType,
		PageSize: limit,
		Offset:   offset,
	}
	return b
}

// AndGroup 添加AND查询条件组
func (b *QueryConditionBuilder) AndGroup() *QueryConditionBuilder {
	return b.AddGroup(sqlconditions.AND)
}

// OrGroup 添加Or查询条件组
func (b *QueryConditionBuilder) OrGroup() *QueryConditionBuilder {
	return b.AddGroup(sqlconditions.OR)
}

// Equals Equals 条件
func (b *QueryConditionBuilder) Equals(key string, value any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.EQUALS, value)
}

// Contains 包含条件, 前后匹配 %value%
func (b *QueryConditionBuilder) Contains(key string, value any) *QueryConditionBuilder {
	return b.Like(key, fmt.Sprintf("%%%v%%", value))
}

// StartsWith 包含条件, 以字符开头 value%
func (b *QueryConditionBuilder) StartsWith(key string, value any) *QueryConditionBuilder {
	return b.Like(key, fmt.Sprintf("%v%%", value))
}

// EndsWith 包含条件, 以字符结尾 %value
func (b *QueryConditionBuilder) EndsWith(key string, value any) *QueryConditionBuilder {
	return b.Like(key, fmt.Sprintf("%%%v", value))
}

// Like Like 条件
func (b *QueryConditionBuilder) Like(key string, value any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.LIKE, value)
}

// Not Not 条件
func (b *QueryConditionBuilder) Not(key string, value any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.NOT, value)
}

// GreaterThan GreaterThan 条件
func (b *QueryConditionBuilder) GreaterThan(key string, value any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.GREATER_THAN, value)
}

// GreaterThanEqual GreaterThanEqual 条件
func (b *QueryConditionBuilder) GreaterThanEqual(key string, value any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.GREATER_THAN_EQUAL, value)
}

// LessThan LessThan 条件
func (b *QueryConditionBuilder) LessThan(key string, value any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.LESS_THAN, value)
}

// LessThanEqual LessThanEqual 条件
func (b *QueryConditionBuilder) LessThanEqual(key string, value any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.LESS_THAN_EQUAL, value)
}

// In In 条件
func (b *QueryConditionBuilder) In(key string, value []any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.IN, value)
}

// Between Between 条件
func (b *QueryConditionBuilder) Between(key string, value ...any) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.BETWEEN, value)
}

// IsNull IsNull 条件
func (b *QueryConditionBuilder) IsNull(key string) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.IS_NULL, nil)
}

// IsNotNull IsNotNull 条件
func (b *QueryConditionBuilder) IsNotNull(key string) *QueryConditionBuilder {
	return b.AddCondition(key, sqlconditions.IS_NOT_NULL, nil)
}

// JoinGroupsWithAnd 设置组之间的关联条件为 "AND"
func (b *QueryConditionBuilder) JoinGroupsWithAnd() *QueryConditionBuilder {
	return b.SetGroupRelation(sqlconditions.AND)
}

// JoinGroupsWithOr 设置组之间的关联条件为 "OR"
func (b *QueryConditionBuilder) JoinGroupsWithOr() *QueryConditionBuilder {
	return b.SetGroupRelation(sqlconditions.OR)
}

// SetGroupRelation 设置组之间的关联条件, 默认为 "AND"
func (b *QueryConditionBuilder) SetGroupRelation(condition string) *QueryConditionBuilder {
	b.condition = strings.ToUpper(condition)
	return b
}

// AddGroupList 追加QueryConditionGroup
func (b *QueryConditionBuilder) AddGroupList(groups []QueryConditionGroup) *QueryConditionBuilder {
	if len(groups) > 0 {
		for _, group := range groups {
			if len(group.Items) == 0 {
				continue
			}
			b.AddGroup(group.Condition)
			for _, item := range group.Items {
				b.AddCondition(item.Key, item.Condition, item.Value)
			}
		}
	}
	return b
}

// AddGroup 添加查询条件组
func (b *QueryConditionBuilder) AddGroup(condition string) *QueryConditionBuilder {
	// 如果 currentGroup 已经有条件项, 则将其加入到 groups 中
	if b.currentGroup != nil && len(b.currentGroup.Items) > 0 {
		b.groups = append(b.groups, *b.currentGroup)
	}

	// 创建新的条件组并赋值给 currentGroup
	b.currentGroup = &QueryConditionGroup{
		Condition: strings.ToUpper(condition),
		Items:     []QueryConditionItem{},
	}

	return b
}

// AddCondition 添加单个查询条件项
func (b *QueryConditionBuilder) AddCondition(key string, condition string, value any) *QueryConditionBuilder {
	// 确保 currentGroup 已初始化
	if b.currentGroup == nil {
		b.currentGroup = &QueryConditionGroup{
			Condition: sqlconditions.AND,
			Items:     []QueryConditionItem{},
		}
	}

	// 验证条件类型并处理
	condition = strings.ToUpper(condition)
	b.currentGroup.Items = append(b.currentGroup.Items, QueryConditionItem{
		Key:       key,
		Condition: condition,
		Value:     value,
	})
	return b
}

// Build 构建最终的 SimpleQueryCondition 对象
func (b *QueryConditionBuilder) Build() QueryCondition {
	// 创建一个新的groups切片来存储最终结果
	finalGroups := make([]QueryConditionGroup, len(b.groups))
	copy(finalGroups, b.groups)

	// 如果当前组存在且有条件项, 将其添加到finalGroups
	if b.currentGroup != nil && len(b.currentGroup.Items) > 0 {
		finalGroups = append(finalGroups, *b.currentGroup)
	}

	return QueryCondition{
		Condition:  b.condition,
		Groups:     finalGroups,
		OrderBy:    b.orderBy,
		Pagination: b.pagination,
		DBType:     b.dbType,
	}
}
