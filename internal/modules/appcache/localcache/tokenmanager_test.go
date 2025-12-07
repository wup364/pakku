// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package localcache

import (
	"fmt"
	"testing"
	"time"
)

// 测试令牌生成器
func TestTokenManager(t *testing.T) {
	count := 100000
	timeStart := time.Now().Nanosecond()
	tokenManager := (&TokenManager{}).Init()
	tokens := make([]string, count)
	// 测试申请令牌&获取令牌
	for i := 0; i < count; i++ {
		if i > 5 && i < 100 {
			tokens[i] = tokenManager.AskToken(i, int64(i))
		} else {
			tokens[i] = tokenManager.AskToken(i, -1)
		}
		// 立刻读取
		_, ok := tokenManager.GetTokenBody(tokens[i])
		if !ok {
			fmt.Println("无法查询申请的令牌")
			t.FailNow()
		}
	}
	tokensNew := tokenManager.ListTokens()
	timeE := (time.Now().Nanosecond() - timeStart)
	fmt.Printf("计划个数: %d, 实际记录个数: %d, 花费时间: %d Nanosecond ", count, len(tokensNew), timeE)
	if len(tokensNew) != count {
		t.FailNow()
	}
	// 测试令牌过期
	for i := 0; i < count; i++ {
		if i > 5 && i < 100 {
			tokenManager.RefreshToken(tokens[i])
		}
		// 立刻读取
		_, ok := tokenManager.GetTokenBody(tokens[i])
		if !ok {
			fmt.Println("令牌提前过期")
			t.FailNow()
		}
	}
	time.Sleep(time.Duration(10) * time.Second)
	tokensNew = tokenManager.ListTokens()
	timeE = (time.Now().Nanosecond() - timeStart) / 1000000
	fmt.Printf("%d Millisecond后, 计划记录个数: %d, 实际记录个数: %d ", timeE, count, len(tokensNew))
	if len(tokensNew) == count {
		fmt.Println("令牌提未过期")
		t.FailNow()
	}
	// 令牌清空
	tokenManager.Clear()
	tokensNew = tokenManager.ListTokens()
	if len(tokensNew) > 0 {
		fmt.Println("令牌无法清空")
		t.FailNow()
	}
	//
	type testStruct struct {
		body1 string
	}
	token := tokenManager.AskToken(&testStruct{body1: "1"}, -1)
	tokenbody, _ := tokenManager.GetTokenBody(token)
	test := tokenbody.(*testStruct)
	fmt.Println(test)
	test.body1 = "2"
	tokenbody, _ = tokenManager.GetTokenBody(token)
	test = tokenbody.(*testStruct)
	fmt.Println(test.body1)
	//
	tokenManager.Destroy()
}

// 令牌生成器
func BenchmarkTokenManager(t *testing.B) {
	count := 5
	timeStart := time.Now().Nanosecond()
	tokenManager := (&TokenManager{}).Init()
	tokens := make([]string, count)
	// 测试申请令牌&获取令牌
	for i := 0; i < count; i++ {
		tokens[i] = tokenManager.AskToken(i, int64(i))
	}
	// list
	tokensNew := tokenManager.ListTokens()
	timeE := (time.Now().Nanosecond() - timeStart) / 1000000
	fmt.Printf("计划个数: %d, 实际记录个数: %d, 花费时间: %d Millisecond ", count, len(tokensNew), timeE)
	// 测试令牌过期
	for i := 0; i < count; i++ {
		tokenManager.RefreshToken(tokens[i])
	}
	tokensNew = tokenManager.ListTokens()
	timeE = (time.Now().Nanosecond() - timeStart) / 1000000
	fmt.Printf("%d Millisecond后, 计划记录个数: %d, 实际记录个数: %d ", timeE, count, len(tokensNew))
	// 令牌清空
	tokenManager.Clear()
	tokensNew = tokenManager.ListTokens()
	if len(tokensNew) > 0 {
		fmt.Println("令牌无法清空")
		t.FailNow()
	}
	//
	tokenManager.Destroy()
}
