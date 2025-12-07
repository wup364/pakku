// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlquerier

import (
	"database/sql"
)

// Scan sql.Rows.Scan
type Scan func(...any) error

// RowScan 行扫
type RowScan[T any] func(scan Scan) (obj T, err error)

// ScanFirstOneAndClose 扫描返回值只有一个的结果, 并自动关闭 sql.Rows
func ScanFirstOneAndClose[T any](rows *sql.Rows) (res T, err error) {
	defer CloseRowsSilence(rows)
	if rows.Next() {
		err = rows.Scan(&res)
	}
	return
}

// ScanFirstRowAndClose 扫描第一列, 并自动关闭 sql.Rows
func ScanFirstRowAndClose[T any](rows *sql.Rows, rowScan RowScan[T]) (res *T, err error) {
	defer CloseRowsSilence(rows)
	if rows.Next() {
		if dto, err := rowScan(rows.Scan); nil == err {
			res = &dto
		} else {
			return nil, err
		}
	}
	return
}

// ScanAndClose 扫描列表, 并自动关闭 sql.Rows
func ScanAndClose[T any](rows *sql.Rows, rowScan RowScan[T]) (res []T, err error) {
	defer CloseRowsSilence(rows)
	for rows.Next() {
		var dto T
		if dto, err = rowScan(rows.Scan); nil != err {
			res = nil
			return
		}
		res = append(res, dto)
	}
	return
}

// CloseRowsSilence 关闭 sql.Rows, 不反悔error
func CloseRowsSilence(rows *sql.Rows) {
	if nil != rows {
		rows.Close()
	}
}
