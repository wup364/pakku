// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

// 简单查询条件组合sql生成器
package sqlcondition

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/wup364/pakku/pkg/constants/sqlconditions"
	"github.com/wup364/pakku/pkg/logs"
	"github.com/wup364/pakku/pkg/strutil"
)

// NewPageNumberPagination 设置页码模式分页
func NewPageNumberPagination(pageSize, pageNumber int) *Pagination {
	return &Pagination{
		Type:       PageNumberType,
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}
}

// NewLimitOffsetPagination 设置限制偏移模式分页
func NewLimitOffsetPagination(limit, offset int) *Pagination {
	return &Pagination{
		Type:     LimitOffsetType,
		PageSize: limit,
		Offset:   offset,
	}
}

// QueryCondition 简单查询条件
type QueryCondition struct {
	Condition  string                `json:"condition" remark:"查询项目关联关系"`
	Groups     []QueryConditionGroup `json:"groups" remark:"查询条件组"`
	OrderBy    []OrderBy             `json:"orderBy" remark:"排序配置"`
	Pagination *Pagination           `json:"pagination" remark:"分页配置"`
	DBType     string                `json:"dbType" remark:"数据库类型"`
}

// QueryConditionGroup 查询条件组
type QueryConditionGroup struct {
	Condition string               `json:"condition" remark:"查询项目关联关系"`
	Items     []QueryConditionItem `json:"items" remark:"查询条件项"`
}

// QueryConditionItem 查询条件项
type QueryConditionItem struct {
	Key       string `json:"key" remark:"查询提条件的KEY"`
	Condition string `json:"condition" remark:"查询条件"`
	Value     any    `json:"value" remark:"查询提条件的值"`
}

// ConditionSqlAndParams 条件sql语句和参数列表
type ConditionSqlAndParams struct {
	Params           []any  // sql条件参数列表
	BaseConditionSql string // 基本条件参数sql, 不包含分页、排序
	ConditionSql     string // 最终sql 基本sql + 排序sql + 分页sql
	OrderBySQL       string // 排序sql
}

// OrderBy 排序配置
type OrderBy struct {
	Field string `json:"field" remark:"排序字段"`
	Order string `json:"order" remark:"排序方式 ASC/DESC"`
}

// Pagination 分页配置
type Pagination struct {
	Type       PaginationType `json:"type" remark:"分页类型"`
	PageSize   int            `json:"pageSize" remark:"每页数量/限制数量"`
	PageNumber int            `json:"pageNumber" remark:"页码, 从1开始"`
	Offset     int            `json:"offset" remark:"偏移量"`
}

// SetDBType 设置数据库类型
func (sqc QueryCondition) SetDBType(dbType string) QueryCondition {
	sqc.DBType = dbType
	return sqc
}

// AppendGroups 追加QueryConditionGroup
func (sqc QueryCondition) AppendGroups(groups []QueryConditionGroup) QueryCondition {
	if len(groups) > 0 {
		sqc.Groups = append(sqc.Groups, groups...)
	}
	return sqc
}

// GetQuerySql 获取拼接好(baseSql + where + 条件)的SQL查询条件部分语句和对应的参数列表
func (sqc QueryCondition) GetQuerySql(baseSql string) (querySql string, qryArgs []any, err error) {
	var conditionSql ConditionSqlAndParams
	if conditionSql, err = sqc.GetConditionSqlAndParams(); nil != err {
		logs.Error(err)
		return
	}

	querySql = baseSql
	if len(conditionSql.BaseConditionSql) > 0 {
		querySql += " " + sqlconditions.WHERE
	}

	if len(conditionSql.ConditionSql) > 0 {
		querySql += " " + conditionSql.ConditionSql
	}

	logs.Debugf("GetQuerySql \r\nsql: %s \r\nparams: %v ", querySql, conditionSql.Params)
	return querySql, conditionSql.Params, err
}

// GetConditionSqlAndParams SQL查询条件部分语句和对应的参数列表
func (sqc QueryCondition) GetConditionSqlAndParams() (res ConditionSqlAndParams, err error) {
	if err = sqc.VerifyParameters(); err != nil {
		return
	}

	// 构建基础查询条件
	if res.BaseConditionSql, res.Params, err = sqc.GetBaseConditionSqlAndParams(); err != nil {
		return
	}

	// 添加排序
	if res.OrderBySQL, err = sqc.GetOrderBySql(); err != nil {
		return
	}

	// 添加分页
	res.ConditionSql = strings.TrimSpace(res.BaseConditionSql + " " + res.OrderBySQL)
	if res.ConditionSql, err = sqc.buildPaginationClause(res.ConditionSql); err != nil {
		return
	}
	return
}

// GetBaseConditionSqlAndParams 基础查询条件
func (sqc QueryCondition) GetBaseConditionSqlAndParams() (string, []any, error) {
	var params []any
	var groupSQLs []string

	for _, group := range sqc.Groups {
		if len(group.Items) == 0 {
			continue
		}

		groupSQL, groupParams, err := sqc.buildGroupCondition(group)
		if err != nil {
			return "", nil, err
		}

		groupSQLs = append(groupSQLs, groupSQL)
		params = append(params, groupParams...)
	}

	if len(groupSQLs) == 0 {
		return "", nil, nil
	}

	return strings.Join(groupSQLs, " "+sqc.Condition+" "), params, nil
}

// buildGroupCondition 构建单个组的查询条件
func (sqc QueryCondition) buildGroupCondition(group QueryConditionGroup) (string, []any, error) {
	var itemsSQLs []string
	var itemsParams []any

	for _, item := range group.Items {
		itemSQL, itemParams, err := sqc.buildConditionSQL(item)
		if err != nil {
			return "", nil, err
		}
		itemsSQLs = append(itemsSQLs, itemSQL)
		itemsParams = append(itemsParams, itemParams...)
	}

	var groupSQL string
	if len(itemsSQLs) > 1 {
		groupSQL = "(" + strings.Join(itemsSQLs, " "+group.Condition+" ") + ")"
	} else {
		groupSQL = itemsSQLs[0]
	}

	return groupSQL, itemsParams, nil
}

// GetOrderBySql 排序子句
func (sqc QueryCondition) GetOrderBySql() (string, error) {
	if len(sqc.OrderBy) == 0 {
		return "", nil
	}

	var orderByClauses []string
	for _, order := range sqc.OrderBy {
		if !sqc.isValidKey(order.Field) {
			return "", fmt.Errorf("illegal format order by field '%s'", order.Field)
		}
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", order.Field, order.Order))
	}

	return "ORDER BY " + strings.Join(orderByClauses, ", "), nil
}

// buildPaginationClause 构建分页子句
func (sqc QueryCondition) buildPaginationClause(baseSQL string) (string, error) {
	if sqc.Pagination == nil {
		return baseSQL, nil
	}

	if len(sqc.DBType) == 0 {
		return baseSQL, fmt.Errorf("database type is required for pagination")
	}

	return BuildPaginationSql(
		baseSQL,
		sqc.DBType,
		sqc.Pagination.GetLimit(),
		sqc.Pagination.GetOffset(),
	)
}

// buildConditionSQL 处理每个查询条件
func (sqc QueryCondition) buildConditionSQL(item QueryConditionItem) (string, []any, error) {
	// 将条件统一转为大写
	item.Condition = strings.ToUpper(item.Condition)

	switch item.Condition {
	case sqlconditions.IN:
		return sqc.handleInCondition(item)
	case sqlconditions.BETWEEN:
		return sqc.handleBetweenCondition(item)
	case sqlconditions.IS_NULL, sqlconditions.IS_NOT_NULL:
		return sqc.handleNullCondition(item)
	// case sqlconditions.EXISTS, sqlconditions.NOT_EXISTS:
	// 	return sqc.handleExistsCondition(item)
	default:
		return fmt.Sprintf("%s %s ?", item.Key, item.Condition), []any{item.Value}, nil
	}
}

// handleInCondition 处理 IN 查询条件
func (sqc QueryCondition) handleInCondition(item QueryConditionItem) (string, []any, error) {
	values, ok := item.Value.([]any)
	if !ok || len(values) == 0 {
		return "", nil, fmt.Errorf("the value of IN query must be a non empty array")
	}

	// 构建 IN 条件的 SQL 语句
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	inSQL := fmt.Sprintf("%s IN (%s)", item.Key, strings.Join(placeholders, ", "))
	return inSQL, values, nil
}

// handleBetweenCondition 处理 BETWEEN 查询条件
func (sqc QueryCondition) handleBetweenCondition(item QueryConditionItem) (string, []any, error) {
	values, ok := item.Value.([]any)
	if !ok || len(values) != 2 {
		return "", nil, fmt.Errorf("the value queried by BETWEN must be an array of length 2")
	}

	// 构建 BETWEEN 条件的 SQL 语句
	betweenSQL := fmt.Sprintf("%s BETWEEN ? AND ?", item.Key)
	return betweenSQL, []any{values[0], values[1]}, nil
}

// handleNullCondition 处理 IS NULL 或 IS NOT NULL 查询条件
func (sqc QueryCondition) handleNullCondition(item QueryConditionItem) (string, []any, error) {
	if item.Condition != sqlconditions.IS_NULL && item.Condition != sqlconditions.IS_NOT_NULL {
		return "", nil, fmt.Errorf("only supports' IS NULL 'or' IS NOT NULL 'conditions")
	}

	// 处理 IS NULL 或 IS NOT NULL 查询
	nullSQL := fmt.Sprintf("%s %s", item.Key, item.Condition)
	return nullSQL, nil, nil
}

// handleExistsCondition 处理 EXISTS 或 NOT EXISTS 查询条件
// func (sqc SimpleQueryCondition) handleExistsCondition(item QueryConditionItem) (string, []any, error) {
// 	// EXISTS 和 NOT EXISTS 通常涉及子查询, 因此我们需要传递一个查询
// 	subQuery, ok := item.Value.(string)
// 	if !ok || len(subQuery) == 0 {
// 		return "", nil, fmt.Errorf("the EXISTS condition must provide a subquery")
// 	}

// 	// 构建 EXISTS 或 NOT EXISTS 的 SQL 语句
// 	existsSQL := fmt.Sprintf("%s (%s)", item.Condition, subQuery)
// 	return existsSQL, nil, nil
// }

// VerifyParameters 校验参数
func (sqc QueryCondition) VerifyParameters() error {
	if len(sqc.Groups) == 0 {
		return nil // 空的, 不需要验证
	}

	// 将主查询条件转为大写
	sqc.Condition = strings.ToUpper(sqc.Condition)

	if !strutil.EqualsAnyIgnoreCase(sqc.Condition, sqcSupportedGroupConditions...) {
		return fmt.Errorf(sqcUnsupportedConditionErr, sqc.Condition)
	}

	// 循环组验证
	for i := range sqc.Groups {
		group := sqc.Groups[i]
		group.Condition = strings.ToUpper(group.Condition) // 将组条件转为大写
		if !strutil.EqualsAnyIgnoreCase(group.Condition, sqcSupportedGroupConditions...) {
			return fmt.Errorf(sqcUnsupportedConditionErr, group.Condition)
		}

		// 校验当前组的项
		if err := sqc.validGroup(group); err != nil {
			return err
		}
	}

	return nil
}

// validGroup 校验组信息
func (sqc QueryCondition) validGroup(group QueryConditionGroup) error {
	for _, item := range group.Items {
		// 将条件转为大写
		item.Condition = strings.ToUpper(item.Condition)

		if _, supported := sqcSupportedItemConditions[item.Condition]; !supported {
			return fmt.Errorf(sqcUnsupportedConditionErr, item.Condition)
		}
		if len(item.Key) == 0 {
			return errors.New("query key cannot be empty")
		}
		if !sqc.isValidKey(item.Key) {
			return fmt.Errorf("illegal format query key '%s'", item.Key)
		}
	}
	return nil
}

// isValidKey 校验key
func (sqc QueryCondition) isValidKey(key string) bool {
	validKey := regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_.]*[a-zA-Z0-9_]$")
	return validKey.MatchString(key)
}

// GetLimit 获取限制数量
func (p *Pagination) GetLimit() int {
	return p.PageSize
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
	if p.Type == PageNumberType {
		return (p.PageNumber - 1) * p.PageSize
	}
	return p.Offset
}
