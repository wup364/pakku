// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 文件工具

package logs

import (
	"fmt"
	"log"
)

// InfoLogger InfoLogger
func InfoLogger() *log.Logger {
	return logI
}

// SetInfoPrefix SetInfoPrefix
func SetInfoPrefix(prefix string) {
	logI.SetPrefix(prefix)
}

// Info Info
func Info(v ...any) {
	if loggerLeve >= LOG_LEVEL_INFO {
		logI.Output(2, fmt.Sprintln(v...))
	}
}

// Infof Infof
func Infof(format string, v ...any) {
	if loggerLeve >= LOG_LEVEL_INFO {
		logI.Output(2, fmt.Sprintf(format+"\n", v...))
	}
}
