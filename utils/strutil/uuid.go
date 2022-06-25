// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// UUID工具

package strutil

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"time"
)

// 初始化机器ID, 避免每次都算
var machineByte []byte
var lastNano uint64
var lastSeq uint

// 初始化机器ID, 避免每次都算
func init() {
	machineid, err := GetMachineID()
	if nil != err {
		panic(err)
	}
	machineByte, err = hex.DecodeString(machineid)
	if nil != err {
		panic(err)
	}
}

// GetUUID 获取唯一ID(机器+纳秒时间+纳秒内自增序列+随机数)
func GetUUID() string {
	// 当前时间(纳秒) 64bit
	ctime := make([]byte, 8)
	if nownano := uint64(time.Now().UnixNano()); nownano == lastNano {
		// 同一纳秒内, 随机序列数值+1
		lastSeq++
		binary.BigEndian.PutUint64(ctime, nownano)
	} else {
		// 每纳秒形成一个新的随机数值
		lastNano = nownano
		lastSeq = uint(rand.Int31())
		binary.BigEndian.PutUint64(ctime, nownano)
	}
	// 随机数 32bit
	random := make([]byte, 4)
	binary.BigEndian.PutUint32(random, uint32(rand.Int31()))
	// 纳秒内自增序列 32bit
	randomSeq := make([]byte, 4)
	binary.BigEndian.PutUint32(randomSeq, uint32(lastSeq))
	// 汇总计算
	buffer := make([]byte, 0, 16)
	buffer = append(buffer, machineByte[14:]...) // 2
	buffer = append(buffer, ctime...)            // 8
	buffer = append(buffer, randomSeq[1:]...)    // 3
	buffer = append(buffer, random[1:]...)       // 3
	// hex一下
	return hex.EncodeToString(buffer)
}
