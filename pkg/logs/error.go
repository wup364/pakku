// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 文件工具

package logs

import (
	"fmt"
	"log"
	"runtime"
)

// ErrorLogger ErrorLogger
func ErrorLogger() *log.Logger {
	return logE
}

// SetErrorPrefix SetErrorPrefix
func SetErrorPrefix(prefix string) {
	logE.SetPrefix(prefix)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...any) {
	s := fmt.Sprintf(format+"\n", v...)
	if loggerLeve >= LOG_LEVEL_ERROR {
		logE.Output(2, s)
	}
	panic(s)
}

// Panic is equivalent to Println() followed by a call to panic().
func Panic(v ...any) {
	s := fmt.Sprintln(v...)
	if loggerLeve >= LOG_LEVEL_ERROR {
		logE.Output(2, s)
	}
	panic(s)
}

// Error Error
func Error(v ...any) {
	if loggerLeve >= LOG_LEVEL_ERROR {
		logE.Output(2, fmt.Sprintln(v...))
		logWithStackInfo(2)
	}
}

// Errorf Errorf
func Errorf(format string, v ...any) {
	if loggerLeve >= LOG_LEVEL_ERROR {
		logE.Output(2, fmt.Sprintf(format+"\n", v...))
		logWithStackInfo(2)
	}
}

// logWithStackInfo 打印调用栈信息
func logWithStackInfo(startStackLevel int) {
	for i := startStackLevel; ; i++ {
		if pc, file, line, ok := runtime.Caller(i); ok {
			fName := runtime.FuncForPC(pc).Name()
			logEStack.Output(2, fmt.Sprintf("%s\n", fName))
			logEStack.Output(2, fmt.Sprintf("at \t%s:%d\n", file, line))
		} else {
			break
		}
	}
}
