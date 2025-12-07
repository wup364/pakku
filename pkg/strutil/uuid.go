// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// UUID工具

package strutil

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

const (
	// 时钟回拨最大容忍时间
	maxBackwardsDrift = time.Second * 10
	// 序列号最大值
	maxSequence = 0xFFFF
)

var (
	randMu      sync.Mutex
	lastNano    uint64       // 上次的时间戳
	lastSeq     uint16       // 序列号
	timeOffset  uint64       // 时间偏移量
	machineByte []byte       // 机器ID
	rnd         *rand.Rand   // 随机数生成器
	startTime   = time.Now() // 程序启动时间

	// 1. 使用对象池减少内存分配
	bufferPool = sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}

	// 2. 预分配固定大小的字节数组
	timeBytes   = make([]byte, 8)
	seqBytes    = make([]byte, 2)
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
			randMu.Lock()
			rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
			randMu.Unlock()
		}
	}()
}

// GetUUID 生成UUID [机器ID 2字节][时间戳 8字节][序列号 3字节][随机数 3字节]
func GetUUID() string {
	randMu.Lock()
	defer randMu.Unlock()

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

	// 写入序列号 (2字节)
	binary.BigEndian.PutUint16(seqBytes, lastSeq)
	buf.Write(seqBytes[:])

	// 写入随机数 (4字节)
	binary.BigEndian.PutUint32(randomBytes, uint32(rnd.Int31()))
	buf.Write(randomBytes[:])

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
