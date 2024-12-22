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

// ErrTxClosed 事务已经提交
var ErrTxClosed = errors.New("transaction already committed or rolled back")

type anys = []any

// DataSourceInfo 数据源信息
type DataSourceInfo interface {
	// GetDriverName 驱动名字
	GetDriverName() string
}

// SqlExecutor SQL查询器 + SQL执行器
type SqlExecutor interface {
	DataSourceInfo
	Query
	Exec
}

// Query SQL查询器
type Query interface {
	DataSourceInfo

	// Query 查询SQL
	Query(query string, args ...any) (r *sql.Rows, err error)

	// QueryWithPrepare 使用Prepare的方式查询SQL
	QueryWithPrepare(query string, args ...any) (r *sql.Rows, err error)
}

// Exec SQL执行器
type Exec interface {
	DataSourceInfo

	// ExecWith 执行SQL
	Exec(query string, args ...any) (r sql.Result, err error)

	// ExecWithBatch 使用Prepare的方式批量执行SQL
	ExecWithBatch(query string, args ...anys) (r []sql.Result, err error)

	// ExecWithPrepare 使用Prepare的方式执行SQL
	ExecWithPrepare(query string, args ...any) (r sql.Result, err error)
}

// SqlTxExecutor SQL执行器(在事务中)
type SqlTxExecutor interface {
	SqlExecutor

	// Complete 提交事务, 如果出错则使用 RollbackSilence() 回滚事务
	Complete() error

	// Commit 提交事务
	Commit() error

	// Rollback 回滚事务
	Rollback() error

	// RollbackSilence 回滚事务, 不返回错误, 仅记录日志
	RollbackSilence()
}
