// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 拓展对象

package utypes

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSafeMap(t *testing.T) {
	am := SafeMap[int, string]{}
	bm := am.New()
	cm := am.New()
	dm := cm.New()
	for i := 0; i < 10; i++ {
		cm.Put(i, "am-val_"+strconv.Itoa(i))
		bm.Put(i, "bm-val_"+strconv.Itoa(i))
		dm.Put(i, "dm-val_"+strconv.Itoa(i))
	}
	fmt.Println("bm", bm.Size())
	fmt.Println("cm", cm.Size())
	fmt.Println("dm", dm.Size())
	fmt.Println("bm", bm.Keys())
	fmt.Println("bm", bm.Keys())
	fmt.Println("cm", cm.Keys())
	fmt.Println("dm", dm.Keys())
	fmt.Println("bm", bm.Values())
	fmt.Println("cm", cm.Values())
	fmt.Println("dm", dm.Values())
	fmt.Println("bm", bm.ToMap())
	fmt.Println("cm", cm.ToMap())
	fmt.Println("dm", dm.ToMap())
	cm.Clear()
	fmt.Println("bm", bm.Size())
	fmt.Println("cm", cm.Size())
	fmt.Println("dm", dm.Size())
	fmt.Println(bm.CutR())
	fmt.Println(cm.CutR())
	fmt.Println(dm.CutR())
}
