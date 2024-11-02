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
	"sync"

	"github.com/wup364/pakku/utils/logs"
)

// NewSqlExecutor4Tx 带事务的执行器
func NewSqlExecutor4Tx(driverName string, tx *sql.Tx) SqlTxExecutor {
	return &SqlExecutor4Tx{
		driverName: driverName,
		tx:         tx,
	}
}

// SqlExecutor4Tx 在事务中执行
type SqlExecutor4Tx struct {
	driverName string
	committed  bool
	tx         *sql.Tx
	mutex      sync.Mutex
}

// GetDriverName 驱动名字
func (se *SqlExecutor4Tx) GetDriverName() string {
	return se.driverName
}

// GetTx 获取sql.Tx对象
func (se *SqlExecutor4Tx) GetTx() *sql.Tx {
	return se.tx
}

// ExecWith 执行SQL
func (se *SqlExecutor4Tx) Exec(query string, args ...any) (r sql.Result, err error) {
	return se.tx.Exec(query, args...)
}

// ExecWithBatch 使用Prepare的方式批量执行SQL
func (se *SqlExecutor4Tx) ExecWithBatch(query string, args ...anys) (rs []sql.Result, err error) {
	if len(args) == 0 {
		return nil, errors.New("args is empty")
	}

	var stmt *sql.Stmt
	if stmt, err = se.tx.Prepare(query); err != nil {
		return
	}

	for i := 0; i < len(args); i++ {
		var r sql.Result
		if r, err = stmt.Exec(args[i]...); err != nil {
			return
		}
		rs = append(rs, r)
	}
	return
}

// ExecWithPrepare 使用Prepare的方式执行SQL
func (se *SqlExecutor4Tx) ExecWithPrepare(query string, args ...any) (r sql.Result, err error) {
	var stmt *sql.Stmt
	if stmt, err = se.tx.Prepare(query); err != nil {
		return
	}

	return stmt.Exec(args...)
}

// Query 查询SQL
func (se *SqlExecutor4Tx) Query(query string, args ...any) (r *sql.Rows, err error) {
	return se.tx.Query(query, args...)
}

// QueryWithPrepare 使用Prepare的方式查询SQL
func (se *SqlExecutor4Tx) QueryWithPrepare(query string, args ...any) (r *sql.Rows, err error) {
	var stmt *sql.Stmt
	if stmt, err = se.tx.Prepare(query); err != nil {
		return
	}

	return stmt.Query(args...)
}

// Complete 提交事务, 如果出错则使用 RollbackSilence() 回滚事务
func (se *SqlExecutor4Tx) Complete() (err error) {
	if err = se.Commit(); nil != err && !se.committed {
		se.RollbackSilence()
	} else {
		se.tx = nil
	}
	return
}

// RollbackSilence 回滚事务, 不返回错误, 仅记录日志
func (se *SqlExecutor4Tx) RollbackSilence() {
	if err := se.Rollback(); nil != err && !se.committed {
		logs.Errorln(err)
	} else {
		se.tx = nil
	}
}

// Commit 提交事务
func (se *SqlExecutor4Tx) Commit() (err error) {
	se.mutex.Lock()
	defer se.mutex.Unlock()

	if se.committed || nil == se.tx {
		return ErrTxClosed
	}

	if err = se.tx.Commit(); err == nil {
		se.committed = true
	}
	return
}

// Rollback 回滚事务
func (se *SqlExecutor4Tx) Rollback() error {
	se.mutex.Lock()
	defer se.mutex.Unlock()

	if se.committed || nil == se.tx {
		return ErrTxClosed
	}
	return se.tx.Rollback()
}
