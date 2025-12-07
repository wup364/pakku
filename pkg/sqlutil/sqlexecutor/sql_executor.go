// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlexecutor

import (
	"database/sql"
	"errors"
)

// ErrTxClosed 事务已经提交
var ErrTxClosed = errors.New("transaction already committed or rolled back")

type anys = []any

// SqlExecutorProvider sql执行器
type SqlExecutorProvider interface {
	DataSourceInfo

	GetDB() *sql.DB

	// GetSqlExecutor 常规执行器
	GetSqlExecutor() SqlExecutor

	// GetSqlTxExecutor 获取事务执行器(开启一个事务)
	GetSqlTxExecutor() (SqlTxExecutor, error)
}

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
