// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

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

// SortByName 根据名字排序
func SortByName(array []string, asc bool) {
	sort.Sort(&SorterByName{
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

// SorterByName 根据名字排序
type SorterByName struct {
	array []string
	split string
	asc   bool
}

// 实现sort.Interface接口取元素数量方法
func (sort *SorterByName) Len() int {
	return len(sort.array)
}

// 实现sort.Interface接口比较元素方法
func (sort *SorterByName) Less(i, j int) bool {
	less := strings.Compare(sort.array[i], sort.array[j]) < 0
	if !sort.asc {
		less = !less
	}
	return less
}

// 实现sort.Interface接口交换元素方法
func (sort *SorterByName) Swap(i, j int) {
	sort.array[i], sort.array[j] = sort.array[j], sort.array[i]
}
