// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 拓展对象-interface{}转各种类型

package utypes

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSafeMap(t *testing.T) {
	am := SafeMap{}
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
