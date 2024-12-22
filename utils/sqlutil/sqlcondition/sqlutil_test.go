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
	"testing"
)

func TestBuildPaginationSql(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		driverName string
		limit      int
		offset     int
		want       string
		wantErr    bool
	}{
		{
			name:       "MySQL分页测试",
			sql:        "SELECT * FROM users",
			driverName: "mysql",
			limit:      10,
			offset:     0,
			want:       "SELECT * FROM users LIMIT 10 OFFSET 0",
			wantErr:    false,
		},
		{
			name:       "SQLite分页测试",
			sql:        "SELECT * FROM users",
			driverName: "sqlite3",
			limit:      10,
			offset:     20,
			want:       "SELECT * FROM users LIMIT 10 OFFSET 20",
			wantErr:    false,
		},
		{
			name:       "Oracle分页测试",
			sql:        "SELECT * FROM users",
			driverName: "oracle",
			limit:      10,
			offset:     20,
			want:       "SELECT * FROM (SELECT a.*, ROWNUM rn FROM (SELECT * FROM users) a WHERE ROWNUM <= 30) WHERE rn > 20",
			wantErr:    false,
		},
		{
			name:       "MSSQL分页测试",
			sql:        "SELECT * FROM users",
			driverName: "mssql",
			limit:      10,
			offset:     20,
			want:       "SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) AS RowNum, t.* FROM (SELECT * FROM users) t) AS temp WHERE RowNum BETWEEN 21 AND 30",
			wantErr:    false,
		},
		{
			name:       "不支持的数据库类型测试",
			sql:        "SELECT * FROM users",
			driverName: "unknown",
			limit:      10,
			offset:     0,
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildPaginationSql(tt.sql, tt.driverName, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildPaginationSql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildPaginationSql() = %v, want %v", got, tt.want)
			}
		})
	}
}
