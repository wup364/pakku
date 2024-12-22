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
