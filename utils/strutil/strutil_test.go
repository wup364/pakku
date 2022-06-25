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
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func TestGetUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println("GetUUID", GetUUID())
	}
	for i := 0; i < 10; i++ {
		fmt.Println("GetUUID", GetUUID())
		time.Sleep(time.Nanosecond)
	}
}

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
