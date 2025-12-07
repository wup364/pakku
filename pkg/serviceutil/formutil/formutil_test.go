// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

// 提供 http form 参数转换工具
package formutil

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// 测试数据结构
type TestStruct struct {
	Name     string        `form:"name" default:"defaultName"`
	Age      int           `form:"age" required:"true"`
	Active   bool          `form:"active"`
	Balance  float64       `form:"balance"`
	Tags     []string      `form:"tags"`
	Duration time.Duration `form:"duration"`
}

func TestFormBinder_Bind(t *testing.T) {
	// 创建一个新的 FormBinder
	binder := NewFormBinder()

	// 模拟 HTTP 请求
	req := httptest.NewRequest("POST", "/", strings.NewReader("name=John&age=30&active=true&balance=100.5&tags=go,programming&duration=1h"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 创建待绑定的结构体
	data := &TestStruct{}

	// 绑定参数
	err := binder.Bind(req, data)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// 检查绑定结果
	if data.Name != "John" {
		t.Errorf("expected Name to be 'John', got '%s'", data.Name)
	}
	if data.Age != 30 {
		t.Errorf("expected Age to be 30, got %d", data.Age)
	}
	if data.Active != true {
		t.Errorf("expected Active to be true, got %v", data.Active)
	}
	if data.Balance != 100.5 {
		t.Errorf("expected Balance to be 100.5, got %f", data.Balance)
	}
	expectedTags := []string{"go", "programming"}
	if !reflect.DeepEqual(data.Tags, expectedTags) {
		t.Errorf("expected Tags to be %v, got %v", expectedTags, data.Tags)
	}
	if data.Duration != time.Hour {
		t.Errorf("expected Duration to be 1h, got %v", data.Duration)
	}
}

func TestFormBinder_Bind_DefaultValue(t *testing.T) {
	// 创建一个新的 FormBinder
	binder := NewFormBinder()

	// 模拟 HTTP 请求
	req := httptest.NewRequest("POST", "/", strings.NewReader("age=30"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 创建待绑定的结构体
	data := &TestStruct{}

	// 绑定参数
	err := binder.Bind(req, data)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// 检查默认值
	if data.Name != "defaultName" {
		t.Errorf("expected default Name to be 'defaultName', got '%s'", data.Name)
	}
}

func BenchmarkFormBinder_Bind(b *testing.B) {
	// 创建一个新的 FormBinder
	binder := NewFormBinder()

	// 模拟 HTTP 请求
	req := httptest.NewRequest("POST", "/", strings.NewReader("name=John&age=30&active=true&balance=100.5&tags=go,programming&duration=1h"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for i := 0; i < b.N; i++ {
		// 创建待绑定的结构体
		data := &TestStruct{}

		// 绑定参数
		err := binder.Bind(req, data)
		if err != nil {
			b.Fatalf("Bind failed: %v", err)
		}
	}
}
