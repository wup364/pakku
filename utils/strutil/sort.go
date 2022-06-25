// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 字符排序

package strutil

import (
	"sort"
	"strings"
)

// SortByLen 根据名字长度排序
func SortByLen(array []string, asc bool) {
	sort.Sort(&SorterByLen{
		asc:   asc,
		array: array,
	})
}

// SortBySplitLen 根据名字指定分割符长度排序
func SortBySplitLen(array []string, split string, asc bool) {
	sort.Sort(&SorterBySplitLen{
		asc:   asc,
		array: array,
		split: split,
	})
}

// SorterByLen 根据名字长度排序
type SorterByLen struct {
	array []string
	asc   bool
}

// 实现sort.Interface接口取元素数量方法
func (sort *SorterByLen) Len() int {
	return len(sort.array)
}

// 实现sort.Interface接口比较元素方法
func (sort *SorterByLen) Less(i, j int) bool {
	less := len(sort.array[i]) < len(sort.array[j])
	if !sort.asc {
		less = !less
	}
	return less
}

// 实现sort.Interface接口交换元素方法
func (sort *SorterByLen) Swap(i, j int) {
	sort.array[i], sort.array[j] = sort.array[j], sort.array[i]
}

// SorterBySplitLen 根据名字指定分割符长度排序
type SorterBySplitLen struct {
	array []string
	split string
	asc   bool
}

// 实现sort.Interface接口取元素数量方法
func (sort *SorterBySplitLen) Len() int {
	return len(sort.array)
}

// 实现sort.Interface接口比较元素方法
func (sort *SorterBySplitLen) Less(i, j int) bool {
	less := len(strings.Split(sort.array[i], sort.split)) < len(strings.Split(sort.array[j], sort.split))
	if !sort.asc {
		less = !less
	}
	return less
}

// 实现sort.Interface接口交换元素方法
func (sort *SorterBySplitLen) Swap(i, j int) {
	sort.array[i], sort.array[j] = sort.array[j], sort.array[i]
}
