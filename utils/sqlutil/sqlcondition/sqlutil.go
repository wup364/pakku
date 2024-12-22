// Copyright (C) 2024 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlcondition

import (
	"fmt"
	"strings"
)

const (
	// SQL 分页模板
	mysqlTemplate  = "%s LIMIT %d OFFSET %d"
	mssqlTemplate  = "SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) AS RowNum, t.* FROM (%s) t) AS temp WHERE RowNum BETWEEN %d AND %d"
	oracleTemplate = "SELECT * FROM (SELECT a.*, ROWNUM rn FROM (%s) a WHERE ROWNUM <= %d) WHERE rn > %d"
)

// BuildPaginationSql 根据不同数据库拼接分页参数
func BuildPaginationSql(sql, driverName string, limit, offset int) (string, error) {
	driverName = strings.ToUpper(driverName)

	// 别名处理
	if strings.Contains(driverName, "SQLITE") {
		driverName = "SQLITE"
	} else if strings.Contains(driverName, "MYSQL") {
		driverName = "MYSQL"
	} else if strings.Contains(driverName, "ORACLE") {
		driverName = "ORACLE"
	} else if strings.Contains(driverName, "POSTGRES") {
		driverName = "POSTGRES"
	} else if strings.Contains(driverName, "MSSQL") {
		driverName = "MSSQL"
	}

	switch driverName {
	case "MYSQL", "SQLITE", "POSTGRES":
		return fmt.Sprintf(mysqlTemplate, sql, limit, offset), nil
	case "ORACLE":
		endRow := offset + limit
		return fmt.Sprintf(oracleTemplate, sql, endRow, offset), nil
	case "MSSQL":
		return fmt.Sprintf(mssqlTemplate, sql, offset+1, offset+limit), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", driverName)
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
