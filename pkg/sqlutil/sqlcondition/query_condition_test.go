// SPDX-License-Identifier: MIT
// Copyright (C) 2023 WuPeng <wup364@outlook.com>.

package sqlcondition

import (
	"reflect"
	"testing"
)

func TestQueryConditionBuilder(t *testing.T) {
	// 测试 1: 创建一个简单的查询
	builder := NewQueryConditionBuilder().SetDBType("MYSQL")

	query := builder.
		OrGroup().
		GreaterThan("age", 30).
		Contains("name", "John").
		AndGroup().
		In("status", []any{"active", "pending"}).
		Between("created_at", "2024-01-01", "2024-12-31").
		JoinGroupsWithAnd().
		// SetLimitOffsetPagination(10, 0).
		OrderByAsc("created_at").Build()

	// 期望的查询条件
	expectedQuery := QueryCondition{
		DBType:    "MYSQL",
		Condition: "AND",
		Groups: []QueryConditionGroup{
			{
				Condition: "OR",
				Items: []QueryConditionItem{
					{Key: "age", Condition: ">", Value: 30},
					{Key: "name", Condition: "LIKE", Value: "%John%"},
				},
			},
			{
				Condition: "AND",
				Items: []QueryConditionItem{
					{Key: "status", Condition: "IN", Value: []interface{}{"active", "pending"}},
					{Key: "created_at", Condition: "BETWEEN", Value: []interface{}{"2024-01-01", "2024-12-31"}},
				},
			},
		},
		OrderBy: []OrderBy{
			OrderBy{
				Field: "created_at",
				Order: "ASC",
			},
		},
	}

	// 比较结果与期望
	if !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Expected query: %+v, but got: %+v", expectedQuery, query)
	}

	// 测试 2: 测试没有条件的情况
	builder = NewQueryConditionBuilder().SetDBType("MYSQL")
	query2 := builder.SetGroupRelation("OR").
		AddCondition("name", "=", "Alice").
		Build()

	expectedQuery2 := QueryCondition{
		Condition: "OR",
		DBType:    "MYSQL",
		Groups: []QueryConditionGroup{
			{
				Condition: "AND", // 默认 AND
				Items: []QueryConditionItem{
					{Key: "name", Condition: "=", Value: "Alice"},
				},
			},
		},
	}

	if !reflect.DeepEqual(query2, expectedQuery2) {
		t.Errorf("Expected query: %+v, but got: %+v", expectedQuery2, query2)
	}

	// 测试 3: 测试只有一个条件组的情况
	builder = NewQueryConditionBuilder().SetDBType("MYSQL")
	query3 := builder.AddGroup("AND").
		AddCondition("age", ">", 25).
		Build()

	expectedQuery3 := QueryCondition{
		Condition: "AND",
		DBType:    "MYSQL",
		Groups: []QueryConditionGroup{
			{
				Condition: "AND",
				Items: []QueryConditionItem{
					{Key: "age", Condition: ">", Value: 25},
				},
			},
		},
	}

	if !reflect.DeepEqual(query3, expectedQuery3) {
		t.Errorf("Expected query: %+v, but got: %+v", expectedQuery3, query3)
	}

}

func TestGetSQLConditionAndParams(t *testing.T) {
	// 测试用例 1: 生成包含 AND 和 OR 条件的查询
	builder := NewQueryConditionBuilder().SetDBType("MYSQL")

	// 构建查询条件
	query := builder.
		OrGroup().
		GreaterThan("age", 30).
		Like("name", "%John%").
		AndGroup().
		In("status", []any{"active", "pending"}).
		Between("created_at", "2024-01-01", "2024-12-31").
		SetLimitOffsetPagination(10, 0).
		OrderByAsc("created_at").
		Build()

	// 调用 GetSQLConditionAndParams 获取 SQL 和参数
	conditionSqlAndParams, err := query.GetConditionSqlAndParams()
	if nil != err {
		panic(err)
	}

	// 期望的 SQL 查询条件
	expectedSQL := "(age > ? OR name LIKE ?) AND (status IN (?, ?) AND created_at BETWEEN ? AND ?) ORDER BY created_at ASC LIMIT 10 OFFSET 0"
	expectedParams := []interface{}{30, "%John%", "active", "pending", "2024-01-01", "2024-12-31"}

	// 检查 SQL 查询语句
	if conditionSqlAndParams.ConditionSql != expectedSQL {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL, conditionSqlAndParams.ConditionSql)
	}

	// 检查查询参数
	if !reflect.DeepEqual(conditionSqlAndParams.Params, expectedParams) {
		t.Errorf("Expected Params: %+v, but got: %+v", expectedParams, conditionSqlAndParams.Params)
	}

	// 测试用例 2: 测试只有一个简单条件的情况
	builder = NewQueryConditionBuilder().SetDBType("MYSQL")
	query2 := builder.Equals("name", "Alice").Build()

	conditionSqlAndParams2, err := query2.GetConditionSqlAndParams()
	if nil != err {
		panic(err)
	}

	expectedSQL2 := "name = ?"
	expectedParams2 := []interface{}{"Alice"}

	// 检查 SQL 查询语句
	if conditionSqlAndParams2.BaseConditionSql != expectedSQL2 {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL2, conditionSqlAndParams2.BaseConditionSql)
	}

	// 检查查询参数
	if !reflect.DeepEqual(conditionSqlAndParams2.Params, expectedParams2) {
		t.Errorf("Expected Params: %+v, but got: %+v", expectedParams2, conditionSqlAndParams2.Params)
	}

	// 测试用例 3: 测试 IN 条件
	builder = NewQueryConditionBuilder().SetDBType("MYSQL")
	query3 := builder.In("status", []any{"active", "pending"}).Build()

	conditionSqlAndParams3, err := query3.GetConditionSqlAndParams()
	if nil != err {
		panic(err)
	}

	expectedSQL3 := "status IN (?, ?)"
	expectedParams3 := []interface{}{"active", "pending"}

	// 检查 SQL 查询语句
	if conditionSqlAndParams3.ConditionSql != expectedSQL3 {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL3, conditionSqlAndParams3.ConditionSql)
	}

	// 检查查询参数
	if !reflect.DeepEqual(conditionSqlAndParams3.Params, expectedParams3) {
		t.Errorf("Expected Params: %+v, but got: %+v", expectedParams3, conditionSqlAndParams3.Params)
	}

	// 测试用例 4: 测试 BETWEEN 条件
	builder = NewQueryConditionBuilder().SetDBType("MYSQL")
	query4 := builder.Between("created_at", "2024-01-01", "2024-12-31").Build()

	conditionSqlAndParams4, err := query4.GetConditionSqlAndParams()
	if nil != err {
		panic(err)
	}

	expectedSQL4 := "created_at BETWEEN ? AND ?"
	expectedParams4 := []interface{}{"2024-01-01", "2024-12-31"}

	// 检查 SQL 查询语句
	if conditionSqlAndParams4.ConditionSql != expectedSQL4 {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL4, conditionSqlAndParams4.ConditionSql)
	}

	// 检查查询参数
	if !reflect.DeepEqual(conditionSqlAndParams4.Params, expectedParams4) {
		t.Errorf("Expected Params: %+v, but got: %+v", expectedParams4, conditionSqlAndParams4.Params)
	}

	// 测试用例 5: 测试多个条件组合
	builder = NewQueryConditionBuilder().SetDBType("MYSQL")
	query5 := builder.AndGroup().
		GreaterThan("age", 30).
		Like("name", "%John%").
		OrGroup().
		In("status", []any{"active", "pending"}).
		Between("created_at", "2024-01-01", "2024-12-31").
		Build()

	conditionSqlAndParams5, err := query5.GetConditionSqlAndParams()
	if nil != err {
		panic(err)
	}

	expectedSQL5 := "(age > ? AND name LIKE ?) AND (status IN (?, ?) OR created_at BETWEEN ? AND ?)"
	expectedParams5 := []interface{}{30, "%John%", "active", "pending", "2024-01-01", "2024-12-31"}

	// 检查 SQL 查询语句
	if conditionSqlAndParams5.ConditionSql != expectedSQL5 {
		t.Errorf("Expected SQL: %s, but got: %s", expectedSQL5, conditionSqlAndParams5.ConditionSql)
	}

	// 检查查询参数
	if !reflect.DeepEqual(conditionSqlAndParams5.Params, expectedParams5) {
		t.Errorf("Expected Params: %+v, but got: %+v", expectedParams5, conditionSqlAndParams5.Params)
	}
}
