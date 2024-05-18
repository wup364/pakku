// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package ipakku

import "github.com/wup364/pakku/utils/utypes"

// AppConfig app 配置模块
type AppConfig interface {

	// GetConfig 读取key的value信息, 返回 Object 对象, 里面的值可能是string或者map
	GetConfig(key string) utypes.Object

	// SetConfig 设置值
	SetConfig(key string, value interface{}) error

	// ScanAndAutoConfig 扫描带有@autoconfig标签的字段, 并完成其配置
	ScanAndAutoConfig(ptr interface{}) error

	// ScanAndAutoValue 扫描带有@autovalue标签的字段, 并完成其配置
	ScanAndAutoValue(configPrefix string, ptr interface{}) error
}

// IConfig 配置接口
type IConfig interface {

	// Init 初始化解析器
	Init(appName string) error

	// GetConfig 读取key的value信息, 返回 Object 对象, 里面的值可能是string或者map
	GetConfig(key string) (res utypes.Object)

	// SetConfig 设置值
	SetConfig(key string, value interface{}) error
}
