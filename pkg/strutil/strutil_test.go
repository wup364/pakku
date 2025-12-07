// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// UUID工具

package strutil

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestSortByLen(t *testing.T) {
	array := []string{"/api/user", "/", "/api", "/api/xxxx"}
	SortByLen(array, true)
	fmt.Println(array)
}

func TestParse2UnixPath(t *testing.T) {
	path := "http://./user"
	unixpath := Parse2UnixPath(path)
	fmt.Println(unixpath)
	fmt.Println(filepath.Abs(path))
	fmt.Println(filepath.Abs(unixpath))
}

func TestGetSHA256(t *testing.T) {
	fmt.Println(GetSHA256("DN101@4cbc16e9a1e02bb169b4629a0f104dc7"))
	fmt.Println(GetMD5(GetSHA256("DN101@4cbc16e9a1e02bb169b4629a0f104dc7")))
}

func TestDependencySorter(t *testing.T) {
	moduleA := DS_M{Name: "A", Dependencies: []string{"B", "C"}}
	moduleB := DS_M{Name: "B", Dependencies: []string{"C"}}
	moduleC := DS_M{Name: "C", Dependencies: []string{}}
	moduleD := DS_M{Name: "D", Dependencies: []string{"A", "B"}}
	moduleE := DS_M{Name: "E", Dependencies: []string{"D"}}

	order := DependencySorter(moduleA, moduleB, moduleC, moduleD, moduleE)
	for _, module := range order {
		fmt.Println(module)
	}
}
