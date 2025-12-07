// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 文件工具

package logs

import (
	"fmt"
	"log"
)

// DebugLogger DebugLogger
func DebugLogger() *log.Logger {
	return logD
}

// SetDebugPrefix SetDebugPrefix
func SetDebugPrefix(prefix string) {
	logD.SetPrefix(prefix)
}

// Debug Debug
func Debug(v ...any) {
	if loggerLeve >= LOG_LEVEL_DEBUG {
		logD.Output(2, fmt.Sprintln(v...))
	}
}

// Debugf Debugf
func Debugf(format string, v ...any) {
	if loggerLeve >= LOG_LEVEL_DEBUG {
		logD.Output(2, fmt.Sprintf(format+"\n", v...))
	}
}
