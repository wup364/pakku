// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlutil

import (
	"github.com/wup364/pakku/pkg/sqlutil/sqlcondition"
	"github.com/wup364/pakku/pkg/sqlutil/sqlquerier"
)

// SqlQuerier sql查询器
type SqlQuerier[T any] sqlquerier.SqlQuerier[T]

// Scan sql.Rows.Scan
type Scan = sqlquerier.Scan

// QueryCondition 简单查询条件
type QueryCondition = sqlcondition.QueryCondition

// NewSqlQuerier 条件查询执行器
func NewSqlQuerier[T any](baseSql string, conditions sqlcondition.QueryCondition) SqlQuerier[T] {
	return SqlQuerier[T](sqlquerier.NewSimpleSqlQuerier[T](baseSql, conditions))
}

// NewQueryConditionBuilder 创建一个新的查询构建器
func NewQueryConditionBuilder() *sqlcondition.QueryConditionBuilder {
	return sqlcondition.NewQueryConditionBuilder()
}

// NewPageNumberPagination 设置页码模式分页
func NewPageNumberPagination(pageSize, pageNumber int) *sqlcondition.Pagination {
	return sqlcondition.NewPageNumberPagination(pageSize, pageNumber)
}

// NewLimitOffsetPagination 设置限制偏移模式分页
func NewLimitOffsetPagination(limit, offset int) *sqlcondition.Pagination {
	return sqlcondition.NewLimitOffsetPagination(limit, offset)
}

// BuildPaginationSql 根据不同数据库拼接分页参数
func BuildPaginationSql(sql, driverName string, limit, offset int) (string, error) {
	return sqlcondition.BuildPaginationSql(sql, driverName, limit, offset)
}

// SqlPlaceholdersGenerate SQL查询条件占位符生成
func SqlPlaceholdersGenerate(length int) (res string) {
	return sqlcondition.SqlPlaceholdersGenerate(length)
}
