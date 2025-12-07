// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlutil

import (
	"database/sql"

	"github.com/wup364/pakku/pkg/sqlutil/sqlexecutor"
)

// NewSqlExecutorProvider 获取sql执行器, 可从实例中获取普通无事务执行器和带事务的执行器
func NewSqlExecutorProvider(driverName string, db *sql.DB) sqlexecutor.SqlExecutorProvider {
	return sqlexecutor.NewSimpleSqlExecutorProvider(driverName, db)
}

// NewSqlExecutor 普通无事务执行器
func NewSqlExecutor(driverName string, db *sql.DB) sqlexecutor.SqlExecutor {
	return sqlexecutor.NewSqlExecutor4Normal(driverName, db)
}

// NewSqlExecutor4Tx 带事务的执行器
func NewSqlExecutor4Tx(driverName string, tx *sql.Tx) sqlexecutor.SqlTxExecutor {
	return sqlexecutor.NewSqlExecutor4Tx(driverName, tx)
}
