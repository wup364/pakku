// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlquerier

import (
	"database/sql"

	"github.com/wup364/pakku/pkg/sqlutil/sqlcondition"
	"github.com/wup364/pakku/pkg/sqlutil/sqlexecutor"
)

// NewSimpleSqlQuerier 简单条件查询器
func NewSimpleSqlQuerier[T any](baseSql string, conditions sqlcondition.QueryCondition) SqlQuerier[T] {
	return SimpleSqlQuerier[T]{
		BaseSql:      baseSql,
		SqlCondition: conditions,
	}
}

// SimpleSqlQuerier 简单查询
type SimpleSqlQuerier[T any] struct {
	BaseSql      string                      // 基础sql
	SqlCondition sqlcondition.QueryCondition //  简单查询条件
}

// QueryFirstOne 查询第一行第一列
func (sqc SimpleSqlQuerier[T]) QueryFirstOne(query sqlexecutor.Query) (res T, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(query.GetDriverName())
	if querySql, qryArgs, err := sqc.SqlCondition.GetQuerySql(sqc.BaseSql); nil != err {
		return res, err
	} else {
		return QueryFirstOne[T](query, querySql, qryArgs)
	}
}

// QueryFirstRow 查询第一行
func (sqc SimpleSqlQuerier[T]) QueryFirstRow(query sqlexecutor.Query, scan RowScan[T]) (res *T, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(query.GetDriverName())
	if querySql, qryArgs, err := sqc.SqlCondition.GetQuerySql(sqc.BaseSql); nil != err {
		return res, err
	} else {
		return QueryFirstRow(query, querySql, qryArgs, scan)
	}
}

// QueryList 查询结果列表
func (sqc SimpleSqlQuerier[T]) QueryList(query sqlexecutor.Query, scan RowScan[T]) (res []T, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(query.GetDriverName())
	if querySql, qryArgs, err := sqc.SqlCondition.GetQuerySql(sqc.BaseSql); nil != err {
		return res, err
	} else {
		return QueryList(query, querySql, qryArgs, scan)
	}
}

// GetSql 获得组装后的sql
func (sqc SimpleSqlQuerier[T]) GetSql(driverName string) (querySql string, qryArgs []any, err error) {
	sqc.SqlCondition = sqc.SqlCondition.SetDBType(driverName)
	return sqc.SqlCondition.GetQuerySql(sqc.BaseSql)
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
