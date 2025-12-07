package utypes

import (
	"testing"
	"time"
)

type testStructType struct {
	Key1 string
	Key2 string
	Data string
}

func TestObjectScan(t *testing.T) {
	// 测试基础类型转换
	tests := []struct {
		name     string
		input    any
		target   any
		expected any
		wantErr  bool
	}{
		{
			name:     "scan string",
			input:    "hello",
			target:   new(string),
			expected: "hello",
		},
		{
			name:     "scan []byte",
			input:    []byte("hello"),
			target:   new([]byte),
			expected: []byte("hello"),
		},
		{
			name:     "scan int",
			input:    123,
			target:   new(int),
			expected: 123,
		},
		{
			name:     "scan int",
			input:    "123",
			target:   new(int),
			expected: 123,
		},
		{
			name:     "scan int8",
			input:    int8(8),
			target:   new(int8),
			expected: int8(8),
		},
		{
			name:     "scan int16",
			input:    int16(16),
			target:   new(int16),
			expected: int16(16),
		},
		{
			name:     "scan int32",
			input:    int32(32),
			target:   new(int32),
			expected: int32(32),
		},
		{
			name:     "scan int64",
			input:    int64(64),
			target:   new(int64),
			expected: int64(64),
		},
		{
			name:     "scan uint",
			input:    uint(123),
			target:   new(uint),
			expected: uint(123),
		},
		{
			name:     "scan uint8",
			input:    uint8(8),
			target:   new(uint8),
			expected: uint8(8),
		},
		{
			name:     "scan uint16",
			input:    uint16(16),
			target:   new(uint16),
			expected: uint16(16),
		},
		{
			name:     "scan uint32",
			input:    uint32(32),
			target:   new(uint32),
			expected: uint32(32),
		},
		{
			name:     "scan uint64",
			input:    uint64(64),
			target:   new(uint64),
			expected: uint64(64),
		},
		{
			name:     "scan float32",
			input:    float32(3.14),
			target:   new(float32),
			expected: float32(3.14),
		},
		{
			name:     "scan float64",
			input:    float64(3.14159),
			target:   new(float64),
			expected: float64(3.14159),
		},
		{
			name:     "scan bool",
			input:    true,
			target:   new(bool),
			expected: true,
		},
		{
			name:     "scan time",
			input:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			target:   new(time.Time),
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "scan nil",
			input:    nil,
			target:   new(string),
			expected: "",
		},
		{
			name:     "scan map",
			input:    map[string]string{"a": "b"},
			target:   new(map[string]string),
			expected: map[string]string{"a": "b"},
		},
		{
			name:     "scan struct",
			input:    testStructType{Key1: "Key1", Key2: "Key2"},
			target:   new(testStructType),
			expected: testStructType{Key1: "Key1", Key2: "Key2"},
		},
		{
			name:     "scan struct",
			input:    map[string]string{"Key1": "Key1", "Key2": "Key2"},
			target:   new(testStructType),
			expected: testStructType{Key1: "Key1", Key2: "Key2"},
		},
		{
			name:     "scan map",
			input:    &testStructType{Key1: "Key1", Key2: "Key2"},
			target:   new(map[string]string),
			expected: map[string]string{"Key1": "Key1", "Key2": "Key2", "Data": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := Object{o: tt.input}
			err := obj.Scan(tt.target)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("Object.Scan() name=%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// 根据类型进行值比较
			switch target := tt.target.(type) {
			case *string:
				if *target != tt.expected.(string) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *[]byte:
				if string(*target) != string(tt.expected.([]byte)) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *int:
				if *target != tt.expected.(int) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *int8:
				if *target != tt.expected.(int8) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *int16:
				if *target != tt.expected.(int16) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *int32:
				if *target != tt.expected.(int32) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *int64:
				if *target != tt.expected.(int64) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *uint:
				if *target != tt.expected.(uint) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *uint8:
				if *target != tt.expected.(uint8) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *uint16:
				if *target != tt.expected.(uint16) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *uint32:
				if *target != tt.expected.(uint32) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *uint64:
				if *target != tt.expected.(uint64) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *float32:
				if *target != tt.expected.(float32) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *float64:
				if *target != tt.expected.(float64) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *bool:
				if *target != tt.expected.(bool) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			case *time.Time:
				if !(*target).Equal(tt.expected.(time.Time)) {
					t.Errorf("Object.Scan() = %v, want %v", *target, tt.expected)
				}
			}
		})
	}
}

// TestObjectScanEdgeCases 测试一些边界情况
func TestObjectScanEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		obj     Object
		target  any
		wantErr bool
	}{
		{
			name:    "scan to nil target",
			obj:     Object{o: "test"},
			target:  nil,
			wantErr: true,
		},
		{
			name:    "scan from nil object to non-nil target",
			obj:     Object{o: nil},
			target:  new(string),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.obj.Scan(tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Object.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
