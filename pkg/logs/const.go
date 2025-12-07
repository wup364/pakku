// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 文件工具

package logs

import (
	"io"
	"log"
	"os"
)

// LoggerLeve 日志级别 debug info warn error
type LoggerLeve int

const (
	LOG_LEVEL_NONE LoggerLeve = 1 << iota
	LOG_LEVEL_ERROR
	LOG_LEVEL_WARN
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
)

var (
	loggerLeve = LOG_LEVEL_DEBUG
	logD       = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	logI       = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	logW       = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	logE       = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Llongfile)
	logEStack  = log.New(os.Stderr, "", log.Lmsgprefix)
)

// SetOutput 设置输出-info, warn, debug, error
func SetOutput(w io.Writer) {
	logD.SetOutput(w)
	logE.SetOutput(w)
	logI.SetOutput(w)
	logW.SetOutput(w)
	logEStack.SetOutput(w)
}

// SetLoggerLevel NONE DEBUG INFO WARN ERROR
func SetLoggerLevel(lv LoggerLeve) {
	loggerLeve = lv
}
