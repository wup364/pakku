// Copyright (C) 2024 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlutil

import (
	"database/sql"

	"github.com/wup364/pakku/utils/sqlutil/sqlcondition"
	"github.com/wup364/pakku/utils/sqlutil/sqlexecutor"
)

// NewSimpleQuery 简单条件查询器
func NewSimpleQuery[T any](baseSql string, conditions sqlcondition.QueryCondition) SimpleQuery[T] {
	return SimpleQuery[T]{
		BaseSql:      baseSql,
		SqlCondition: conditions,
	}
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

// QueryWithPrepare 执行查询
func QueryWithPrepare(query sqlexecutor.Query, baseSql string, sqlCondition sqlcondition.QueryCondition) (r *sql.Rows, err error) {
	if querySql, qryArgs, err := sqlCondition.GetQuerySql(baseSql); nil != err {
		return nil, err
	} else {
		return query.QueryWithPrepare(querySql, qryArgs...)
	}
}

// QueryFirstOne 查询结果列表
func QueryFirstOne[T any](query sqlexecutor.Query, querySql string, qryArgs []any) (res T, err error) {
	var rows *sql.Rows
	if rows, err = query.QueryWithPrepare(querySql, qryArgs...); nil != err {
		CloseRowsSilence(rows)
		return
	}

	return ScanFirstOneAndClose[T](rows)
}

// QueryFirstRow 查询结果列表
func QueryFirstRow[T any](query sqlexecutor.Query, querySql string, qryArgs []any, scan RowScan[T]) (res *T, err error) {
	var rows *sql.Rows
	if rows, err = query.QueryWithPrepare(querySql, qryArgs...); nil != err {
		CloseRowsSilence(rows)
		return
	}

	return ScanFirstRowAndClose[T](rows, scan)
}

// QueryList 查询结果列表
func QueryList[T any](query sqlexecutor.Query, querySql string, qryArgs []any, scan RowScan[T]) (res []T, err error) {
	var rows *sql.Rows
	if rows, err = query.QueryWithPrepare(querySql, qryArgs...); nil != err {
		CloseRowsSilence(rows)
		return
	}

	return ScanAndClose[T](rows, scan)
}

// SimpleQuery 简单查询
type SimpleQuery[T any] struct {
	BaseSql      string                      // 基础sql
	SqlCondition sqlcondition.QueryCondition //  简单查询条件
}

// QueryFirstOne 查询第一行第一列
func (sqc SimpleQuery[T]) QueryFirstOne(query sqlexecutor.Query) (res T, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(query.GetDriverName())
	if querySql, qryArgs, err := sqc.SqlCondition.GetQuerySql(sqc.BaseSql); nil != err {
		return res, err
	} else {
		return QueryFirstOne[T](query, querySql, qryArgs)
	}
}

// QueryFirstRow 查询第一行
func (sqc SimpleQuery[T]) QueryFirstRow(query sqlexecutor.Query, scan RowScan[T]) (res *T, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(query.GetDriverName())
	if querySql, qryArgs, err := sqc.SqlCondition.GetQuerySql(sqc.BaseSql); nil != err {
		return res, err
	} else {
		return QueryFirstRow(query, querySql, qryArgs, scan)
	}
}

// QueryList 查询结果列表
func (sqc SimpleQuery[T]) QueryList(query sqlexecutor.Query, scan RowScan[T]) (res []T, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(query.GetDriverName())
	if querySql, qryArgs, err := sqc.SqlCondition.GetQuerySql(sqc.BaseSql); nil != err {
		return res, err
	} else {
		return QueryList(query, querySql, qryArgs, scan)
	}
}

// GetSql 获得组装后的sql
func (sqc SimpleQuery[T]) GetSql(driverName string) (querySql string, qryArgs []any, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(driverName)
	return sqc.SqlCondition.GetQuerySql(sqc.BaseSql)
}
