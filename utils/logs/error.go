// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if loggerLeve >= ERROR {
		logE.Output(2, s)
	}
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	if loggerLeve >= ERROR {
		logE.Output(2, s)
	}
	panic(s)
}

// Error Error
func Error(v ...interface{}) {
	if loggerLeve >= ERROR {
		logE.Output(2, fmt.Sprint(v...))
		logWithStackInfo(2)
	}
}

// Errorf Errorf
func Errorf(format string, v ...interface{}) {
	if loggerLeve >= ERROR {
		logE.Output(2, fmt.Sprintf(format, v...))
		logWithStackInfo(2)
	}
}

// Errorln Errorln
func Errorln(v ...interface{}) {
	if loggerLeve >= ERROR {
		logE.Output(2, fmt.Sprintln(v...))
		logWithStackInfo(2)
	}
}

// logWithStackInfo 打印调用栈信息
func logWithStackInfo(startStackLevel int) {
	if errLogStackLevel <= 0 {
		return
	}
	for i := startStackLevel; i < errLogStackLevel; i++ {
		if pc, file, line, ok := runtime.Caller(i); ok {
			fName := runtime.FuncForPC(pc).Name()
			logEStack.Output(2, fmt.Sprintf("%s\n", fName))
			logEStack.Output(2, fmt.Sprintf("at \t%s:%d\n", file, line))
		}
	}
}
