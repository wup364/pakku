// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package sqlexecutor

import (
	"database/sql"
)

// NewSimpleSqlExecutorProvider 获取sql执行器包裹
func NewSimpleSqlExecutorProvider(driverName string, db *sql.DB) SqlExecutorProvider {
	return &SimpleSqlExecutorProvider{
		db:          db,
		driverName:  driverName,
		sqlExecutor: NewSqlExecutor4Normal(driverName, db),
	}
}

// SimpleSqlExecutorProvider sql执行器包裹
type SimpleSqlExecutorProvider struct {
	db          *sql.DB
	driverName  string
	sqlExecutor SqlExecutor
}

// GetDriverName 驱动名字
func (ssep *SimpleSqlExecutorProvider) GetDriverName() string {
	return ssep.driverName
}

// GetDB 获取数据库连接
func (ssep *SimpleSqlExecutorProvider) GetDB() *sql.DB {
	return ssep.db
}

// GetSqlExecutor 常规执行器
func (ssep *SimpleSqlExecutorProvider) GetSqlExecutor() SqlExecutor {
	return ssep.sqlExecutor
}

// GetSqlTxExecutor 获取事务执行器
func (ssep *SimpleSqlExecutorProvider) GetSqlTxExecutor() (res SqlTxExecutor, err error) {
	var tx *sql.Tx
	if tx, err = ssep.GetDB().Begin(); nil == err {
		res = NewSqlExecutor4Tx(ssep.GetDriverName(), tx)
	}
	return
}
