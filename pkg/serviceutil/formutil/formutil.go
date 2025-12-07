// SPDX-License-Identifier: MIT
// Copyright (C) 2024 WuPeng <wup364@outlook.com>.

// 提供 http form 参数转换工具
package formutil

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FormBinder form 参数绑定器
type FormBinder struct {
	TagName string // tag 名称，默认为 "form"
}

// NewFormBinderWithTagName 创建一个新的 FormBinder，使用默认的 tag 名称 "form"
func NewFormBinder() *FormBinder {
	return NewFormBinderWithTagName("")
}

// NewFormBinder 创建一个新的 FormBinder
func NewFormBinderWithTagName(tagName string) *FormBinder {
	if tagName == "" {
		tagName = "form"
	}
	return &FormBinder{
		TagName: tagName,
	}
}

// validateInput 验证绑定的输入参数
func (b *FormBinder) validateInput(r *http.Request, obj interface{}) (reflect.Value, error) {
	if r == nil {
		return reflect.Value{}, fmt.Errorf("request is nil")
	}

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("obj must be a pointer")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("obj must be a pointer to struct")
	}
	return val, nil
}

// getFieldValue 获取字段的 form 值，必要时应用默认值
func (b *FormBinder) getFieldValue(r *http.Request, formTag, defaultVal string) string {
	formVal := r.FormValue(formTag)
	if formVal == "" {
		formVal = defaultVal
	}
	return formVal
}

// validateRequiredField 检查必填字段是否有值
func (b *FormBinder) validateRequiredField(formTag, value string, required bool) error {
	if value == "" && required {
		return fmt.Errorf("field %s is required", formTag)
	}
	return nil
}

// Bind 将 form 参数绑定到结构体
// 支持的 tag:
//   - form:"name"       参数名
//   - default:"value"   默认值
//   - required:"true"   是否必须
func (b *FormBinder) Bind(r *http.Request, obj interface{}) error {
	value, err := b.validateInput(r, obj)
	if err != nil {
		return err
	}

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		structField := typ.Field(i)

		if !field.CanSet() {
			continue
		}
		// 如果字段是结构体，递归处理
		if field.Kind() == reflect.Struct {
			// 递归绑定嵌入结构体
			fieldAddr := field.Addr()
			if fieldAddr.IsValid() && fieldAddr.CanInterface() {
				if err := b.Bind(r, fieldAddr.Interface()); err != nil {
					return err
				}
			}
			continue
		}

		formTag := structField.Tag.Get(b.TagName)
		if formTag == "" {
			continue
		}

		defaultVal := structField.Tag.Get("default")
		required := structField.Tag.Get("required") == "true"

		formVal := b.getFieldValue(r, formTag, defaultVal)

		if err := b.validateRequiredField(formTag, formVal, required); err != nil {
			return err
		}

		if err := b.setField(field, formVal); err != nil {
			return fmt.Errorf("field %s: %v", formTag, err)
		}
	}

	return nil
}

// setField 设置字段值
func (b *FormBinder) setField(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}
	return b.convertAndSetValue(field, value)
}

// convertAndSetValue 处理字段值的类型转换和设置
func (b *FormBinder) convertAndSetValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Bool:
		return b.setBoolField(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return b.setIntField(field, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return b.setUintField(field, value)
	case reflect.Float32, reflect.Float64:
		return b.setFloatField(field, value)
	case reflect.Slice:
		return b.setSliceField(field, value)
	default:
		return fmt.Errorf("unsupported type: %v", field.Type())
	}
}

// setBoolField 设置布尔字段
func (b *FormBinder) setBoolField(field reflect.Value, value string) error {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	field.SetBool(v)
	return nil
}

// setIntField 设置整数字段
func (b *FormBinder) setIntField(field reflect.Value, value string) error {
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		return b.setDurationField(field, value)
	}
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	field.SetInt(v)
	return nil
}

// setUintField 设置无符号整数字段
func (b *FormBinder) setUintField(field reflect.Value, value string) error {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	field.SetUint(v)
	return nil
}

// setFloatField 设置浮点数字段
func (b *FormBinder) setFloatField(field reflect.Value, value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	field.SetFloat(v)
	return nil
}

// setDurationField 设置 duration 字段
func (b *FormBinder) setDurationField(field reflect.Value, value string) error {
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(d))
	return nil
}

// setSliceField 设置 slice 字段
func (b *FormBinder) setSliceField(field reflect.Value, value string) error {
	if field.Type().Elem().Kind() != reflect.String {
		return fmt.Errorf("unsupported slice type: %v", field.Type())
	}
	values := strings.Split(value, ",")
	slice := reflect.MakeSlice(field.Type(), len(values), len(values))
	for i, v := range values {
		slice.Index(i).SetString(strings.TrimSpace(v))
	}
	field.Set(slice)
	return nil
}
