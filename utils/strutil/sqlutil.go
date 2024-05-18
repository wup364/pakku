// Copyright (C) 2023 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package strutil

import (
	"errors"
	"strconv"
	"strings"
)

// SqlConditionConcatForWhere 根据传入参数在原有sql基础上where条件, 自动跳过空值条件
func SqlConditionConcatForWhere(sql, separator string, conditions []string, fields ...string) string {
	sqlConditionStr := SqlConditionConcat(separator, conditions, fields...)
	if len(sqlConditionStr) == 0 {
		return sql
	}
	return sql + " WHERE " + sqlConditionStr
}

// SqlConditionConcat 根据传入参数拼接Sql条件, 自动跳过空值条件
func SqlConditionConcat(separator string, conditions []string, fields ...string) (res string) {
	lenf := len(fields)
	lenc := len(conditions)
	if lenc == 0 || lenf == 0 || lenc < lenf {
		return
	}

	wheres := make([]string, 0)
	for i := 0; i < lenf; i++ {
		if len(fields[i]) > 0 {
			wheres = append(wheres, strings.TrimSpace(conditions[i]))
		}
	}

	if len(wheres) > 0 {
		res = strings.Join(wheres, " "+strings.TrimSpace(separator)+" ")
	}
	return
}

// SqlConditionConcatForPageable 根据不同数据库拼接分页参数, 不支持的driverName返回错误
func SqlConditionConcatForPageable(sql, driverName string, limit, offset int) (newsql string, err error) {
	driverName = strings.ToUpper(driverName)

	// 别名处理
	if strings.Contains(driverName, "SQLITE") {
		driverName = "SQLITE"
	} else if strings.Contains(driverName, "MYSQL") {
		driverName = "MYSQL"
	} else if strings.Contains(driverName, "ORACLE") {
		driverName = "ORACLE"
	}

	//
	switch driverName {
	case "MYSQL", "SQLITE":
		newsql = sql + " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
		return
	case "ORACLE":
		newsql = "SELECT * FROM (SELECT a.*, ROWNUM r FROM (" + sql + ") a WHERE ROWNUM <= " + strconv.Itoa(offset+limit) + ") WHERE r > " + strconv.Itoa(offset)
		return
	default:
		return "", errors.New("unsupported database connection type")
	}
}

// SqlPlaceholdersGenerate SQL查询条件占位符生成
func SqlPlaceholdersGenerate(length int) (res string) {
	if length > 0 {
		for i := 0; i < length; i++ {
			if i == length-1 {
				res += "?"
			} else {
				res += "?,"
			}
		}
	}
	return
}
