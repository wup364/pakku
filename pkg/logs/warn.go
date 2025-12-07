// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 文件工具

package logs

import (
	"fmt"
	"log"
)

// WarnLogger WarnLogger
func WarnLogger() *log.Logger {
	return logW
}

// SetWarnPrefix SetWarnPrefix
func SetWarnPrefix(prefix string) {
	logW.SetPrefix(prefix)
}

// Warn Warn
func Warn(v ...any) {
	if loggerLeve >= LOG_LEVEL_WARN {
		logW.Output(2, fmt.Sprint(v...))
	}
}

// Warnf Warnf
func Warnf(format string, v ...any) {
	if loggerLeve >= LOG_LEVEL_WARN {
		logW.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warnlnf Warnlnf
func Warnlnf(format string, v ...any) {
	if loggerLeve >= LOG_LEVEL_WARN {
		logW.Output(2, fmt.Sprintf(format+"\n", v...))
	}
}

// Warnln Warnln
func Warnln(v ...any) {
	if loggerLeve >= LOG_LEVEL_WARN {
		logW.Output(2, fmt.Sprintln(v...))
	}
}
