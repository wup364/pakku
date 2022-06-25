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
	"io"
	"log"
	"os"
)

// LoggerLeve 日志级别 debug info error
type LoggerLeve int

const (
	// NONE NONE
	NONE LoggerLeve = 1 << iota
	// ERROR ERROR
	ERROR
	// INFO INFO
	INFO
	// DEBUG DEBUG
	DEBUG
)

var (
	loggerLeve = DEBUG
	logI       = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	logD       = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	logE       = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)
)

// SetOutput 设置输出-info, debug, error
func SetOutput(w io.Writer) {
	logD.SetOutput(w)
	logE.SetOutput(w)
	logI.SetOutput(w)
}

// SetLoggerLevel NONE DEBUG INFO ERROR
func SetLoggerLevel(lv LoggerLeve) {
	loggerLeve = lv
}
