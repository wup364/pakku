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
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	// 时钟回拨最大容忍时间
	maxBackwardsDrift = time.Second * 10
	// 序列号最大值
	maxSequence = math.MaxUint32
)

var (
	mu          sync.Mutex
	lastNano    uint64       // 上次的时间戳
	lastSeq     uint32       // 序列号
	timeOffset  uint64       // 时间偏移量
	machineByte []byte       // 机器ID
	rnd         *rand.Rand   // 随机数生成器
	startTime   = time.Now() // 程序启动时间

	// 1. 使用对象池减少内存分配
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// 2. 预分配固定大小的字节数组
	timeBytes   = make([]byte, 8)
	seqBytes    = make([]byte, 4)
	randomBytes = make([]byte, 4)
)

func init() {
	// 初始化随机数生成器
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

	// 初始化机器ID
	if id, err := GetMachineID(); err != nil {
		panic(err)
	} else if machineByte, err = hex.DecodeString(id); err != nil {
		panic(err)
	}

	// 每小时重置随机数种子
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
			mu.Unlock()
		}
	}()
}

// GetUUID 生成UUID [机器ID 2字节][时间戳 8字节][序列号 3字节][随机数 3字节]
func GetUUID() string {
	mu.Lock()
	defer mu.Unlock()

	// 从对象池获取buffer
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// 预分配空间
	if buf.Cap() < 16 {
		buf.Grow(16)
	}

	// 写入机器ID (2字节)
	buf.Write(machineByte[14:])

	// 写入时间戳 (8字节)
	timestamp := getTimestamp()
	binary.BigEndian.PutUint64(timeBytes, timestamp)
	buf.Write(timeBytes)

	// 写入序列号 (3字节)
	binary.BigEndian.PutUint32(seqBytes, lastSeq)
	buf.Write(seqBytes[1:])

	// 写入随机数 (3字节)
	binary.BigEndian.PutUint32(randomBytes, uint32(rnd.Int31()))
	buf.Write(randomBytes[1:])

	return hex.EncodeToString(buf.Bytes())
}

// getTimestamp 获取时间戳并处理时钟回拨
func getTimestamp() uint64 {
	timestamp := uint64(time.Now().UnixNano())

	// 处理时钟回拨
	if timestamp < lastNano {
		drift := lastNano - timestamp

		if drift > uint64(maxBackwardsDrift) {
			// 回拨超过阈值，使用程序运行时间
			timestamp = uint64(startTime.UnixNano() + int64(time.Since(startTime)))
		} else {
			// 小范围回拨，使用偏移量
			timeOffset += drift
			timestamp += timeOffset
		}
	}

	// 确保时间戳递增
	if timestamp <= lastNano {
		if lastSeq == maxSequence {
			timestamp = lastNano + 1
			lastSeq = 0
		} else {
			lastSeq++
		}
	} else {
		lastSeq = 0
	}

	lastNano = timestamp
	return timestamp
}
