// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 文件工具

package logs

import (
	"testing"
)

func TestLogs(t *testing.T) {
	// fs, _ := fileutil.GetWriter("./logs.test.log")
	// SetOutput(fs)
	// SetLoggerLevel(DEBUG)
	Debug("3", "123", "abc")
	Info("1", "123", "abc")
	Error("5", "123", "abc")
}
