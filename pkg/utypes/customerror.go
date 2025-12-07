// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

package utypes

// NewCustomError 创建自定义错误
func NewCustomError(err error, errorCode string) CustomError {
	return &customErrorImpl{
		error:     err,
		errorCode: errorCode,
	}
}

// CustomError 自定义错误
type CustomError interface {
	error

	// ErrorCode 获取业务异常代码
	ErrorCode() string
}

// customErrorImpl 自定义错误实例
type customErrorImpl struct {
	error
	errorCode string
}

func (cr *customErrorImpl) ErrorCode() string {
	return cr.errorCode
}
