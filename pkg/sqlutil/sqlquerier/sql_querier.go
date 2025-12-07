// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlquerier

import (
	"github.com/wup364/pakku/pkg/sqlutil/sqlexecutor"
)

// SqlQuerier sql查询器
type SqlQuerier[T any] interface {

	// QueryFirstOne 查询第一行第一列
	QueryFirstOne(query sqlexecutor.Query) (res T, err error)

	// QueryFirstRow 查询第一行
	QueryFirstRow(query sqlexecutor.Query, scan RowScan[T]) (res *T, err error)

	// QueryList 查询结果列表
	QueryList(query sqlexecutor.Query, scan RowScan[T]) (res []T, err error)

	// GetSql 获得组装后的sql
	GetSql(driverName string) (querySql string, qryArgs []any, err error)
}
