// Copyright (C) 2024 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlexecutor

import (
	"database/sql"
	"errors"
)

// NewSqlExecutor4Normal 普通无事务执行器
func NewSqlExecutor4Normal(driverName string, db *sql.DB) SqlExecutor {
	return &SqlExecutor4Normal{
		driverName: driverName,
		db:         db,
	}
}

// SqlExecutor4Normal 普通执行器(不在事务中执行)
type SqlExecutor4Normal struct {
	driverName string
	db         *sql.DB
}

// GetDriverName 驱动名字
func (se *SqlExecutor4Normal) GetDriverName() string {
	return se.driverName
}

// Exec 执行SQL
func (se *SqlExecutor4Normal) Exec(query string, args ...any) (r sql.Result, err error) {
	return se.db.Exec(query, args...)
}

// ExecWithPrepare 执行SQL, 使用Prepare的方式
func (se *SqlExecutor4Normal) ExecWithPrepare(query string, args ...any) (r sql.Result, err error) {
	var stmt *sql.Stmt
	if stmt, err = se.db.Prepare(query); err != nil {
		return
	}
	return stmt.Exec(args...)
}

// ExecWithBatch 开启一个事务, 在事务中执行SQL, 使用Prepare的方式批量执行
func (se *SqlExecutor4Normal) ExecWithBatch(query string, args ...anys) (rs []sql.Result, err error) {
	if len(args) == 0 {
		return nil, errors.New("args is empty")
	}

	var tx *sql.Tx
	if tx, err = se.db.Begin(); nil != err {
		return
	}

	var stmt *sql.Stmt
	if stmt, err = tx.Prepare(query); err != nil {
		tx.Rollback()
		return
	}

	for i := 0; i < len(args); i++ {
		var r sql.Result
		if r, err = stmt.Exec(args[i]...); err != nil {
			rs = make([]sql.Result, 0)
			tx.Rollback()
			return
		}
		rs = append(rs, r)
	}

	if err = tx.Commit(); nil != err {
		rs = make([]sql.Result, 0)
		tx.Rollback()
	}
	return
}

// Query 查询SQL
func (se *SqlExecutor4Normal) Query(query string, args ...any) (r *sql.Rows, err error) {
	return se.db.Query(query, args...)
}

// QueryWithPrepare 使用Prepare的方式查询SQL
func (se *SqlExecutor4Normal) QueryWithPrepare(query string, args ...any) (r *sql.Rows, err error) {
	var stmt *sql.Stmt
	if stmt, err = se.db.Prepare(query); err != nil {
		return
	}
	return stmt.Query(args...)
}
