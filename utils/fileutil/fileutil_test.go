// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package fileutil

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/wup364/pakku/utils/strutil"
)

func TestGetSHA256(t *testing.T) {
	f, err := OpenFile("C:\\Users\\wupen\\Downloads\\3")
	if nil != err {
		panic(err)
	}
	sha, err := GetSHA256(f)
	if nil != err {
		panic(err)
	}
	fmt.Println(sha)
}
func TestGetSHA256Multi(t *testing.T) {
	TestGetSHA256(t)
	//
	i := 1
	h := sha256.New()
	buf := make([]byte, 1<<20)
	for i <= 2 {
		r, err := OpenFile("C:\\Users\\wupen\\Downloads\\" + strconv.Itoa(i))
		if nil != err {
			panic(err)
		}
		for {
			n, err := io.ReadFull(r, buf)
			if err == nil || err == io.ErrUnexpectedEOF {
				fmt.Println("-2->" + hex.EncodeToString(h.Sum(nil)))
				if _, err = h.Write(buf[0:n]); err != nil {
					panic(err)
				}
			} else if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		i++
	}
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
	//

}

func TestGetDirList(t *testing.T) {
	path := os.TempDir()
	list, _ := GetDirList(path)
	fmt.Println(path, list)
}
func TestMoveFilesAcrossDisk(t *testing.T) {
	src := "C:\\Users\\wupen\\Desktop\\yaml"
	dst := "D:\\.sys\\test\\yaml"
	err := MoveFilesAcrossDisk(src, dst, false, true, func(src, dst string, err error) error {
		fmt.Println(src, "-->", dst, err)
		return err
	})
	fmt.Println(err)
}
func TestGetModifyTime(t *testing.T) {
	osTemp := os.TempDir()
	path := osTemp + "\\TestGetModifyTime\\" + strutil.GetUUID() + ".txt"
	if !IsExist(strutil.GetPathParent(path)) {
		if err := MkdirAll(strutil.GetPathParent(path)); nil != err {
			panic(err)
		}
	}
	WriteTextFile(path, strutil.GetUUID())
	before, _ := GetModifyTime(path)
	time.Sleep(time.Second)
	WriteTextFile(path, strutil.GetUUID())
	after, _ := GetModifyTime(path)
	fmt.Println(after.Unix() - before.Unix())
	if after.Unix()-before.Unix() <= 0 {
		panic("修改时间读取异常")
	}
}
