// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 字符串处理工具

package strutil

import (
	"container/list"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// ReplaceAll -> strings.Replace
func ReplaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

// Parse2UnixPath 删除路径后面 /, 把\转换为/
func Parse2UnixPath(str string) string {
	if len(str) == 0 {
		return ""
	}
	return path.Clean(strings.Replace(str, "\\", "/", -1))
}

// GetPathParent 截取最后一个'/'前的文字
func GetPathParent(path string) string {
	if path == "/" || path == "\\" {
		return ""
	}
	if i := strings.LastIndex(path, "\\"); i > -1 {
		return path[:i]
	} else if i := strings.LastIndex(path, "/"); i > -1 {
		return path[:i]
	} else {
		return ""
	}
}

// GetPathName 截取最后一个'/'后的文字
func GetPathName(path string) string {
	if path == "/" || path == "\\" {
		return ""
	}
	path = Parse2UnixPath(path)
	return path[strings.LastIndex(path, "/")+1:]
}

// GetPathSuffix 截取最后一个'.'后的文字
func GetPathSuffix(path string) string {
	if index := strings.LastIndex(path, "."); index > -1 {
		return path[index:]
	}
	return ""
}

// ReadAsString 从io.Reader读取文字
func ReadAsString(src io.Reader) string {
	if nil == src {
		return ""
	}
	buf := make([]byte, 0)
	for {
		buftemp := make([]byte, 1024)
		nr, er := src.Read(buftemp)
		if nr > 0 {
			buf = append(buf, buftemp[:nr]...)
		}
		if er != nil {
			if er != io.EOF {
				return ""
			}
			break
		}
	}
	return string(buf)
}

// String2Bool 字符转bool
// true -> [1, t, T, true, TRUE, True]
// false -> [0, f, F, false, FALSE, False]
func String2Bool(str string) bool {
	switch strings.ToLower(str) {
	case "1", "t", "true":
		return true
	case "0", "f", "false":
		return false
	}
	return false
}

// Bool2String bool类型转string
func Bool2String(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// String2Int 十进制数字转换
func String2Int(str string, df int) int {
	if len(str) > 0 {
		if i, err := strconv.Atoi(str); nil == err {
			return i
		}
	}
	return df
}

// Json2Struct json转对象
func Json2Struct(jsonStr string, structPointer any) (err error) {
	if len(jsonStr) == 0 {
		return
	}
	err = json.Unmarshal([]byte(jsonStr), structPointer)
	return
}

// StructToJson 对象转json
func StructToJson(obj any) (res string, err error) {
	var bt []byte
	if bt, err = json.Marshal(obj); nil != err {
		return
	}
	res = string(bt)
	return
}

// ToJsonIgnoreError 对象转json, 忽略错误
func ToJsonIgnoreError(obj any) (res string) {
	res, _ = StructToJson(obj)
	return
}

// EqualsAny 判断字符在一个数组中存在
func EqualsAny(str string, arr ...string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return true
		}
	}
	return false
}

// EqualsAnyIgnoreCase 判断字符在一个数组中存在, 忽略大小写
func EqualsAnyIgnoreCase(str string, arr ...string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str || strings.EqualFold(arr[i], str) {
			return true
		}
	}
	return false
}

// StartsWithAny 判断字符以数组中任意一条数据开头
func StartsWithAny(str string, arr ...string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str || strings.HasPrefix(str, arr[i]) {
			return true
		}
	}
	return false
}

// StartsWithAnyIgnoreCase 判断字符以数组中任意一条数据开头, 忽略大小写
func StartsWithAnyIgnoreCase(str string, arr ...string) bool {
	strLower := strings.ToLower(str)
	for i := 0; i < len(arr); i++ {
		if strings.EqualFold(arr[i], strLower) || strings.HasPrefix(strLower, strings.ToLower(arr[i])) {
			return true
		}
	}
	return false
}

// GetMD5 字符转MD5
func GetMD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

// GetSHA256 字符转sha256
func GetSHA256(str string) string {
	hx := sha256.New()
	hx.Write([]byte(str))
	return hex.EncodeToString(hx.Sum(nil))
}

// GetMachineID 放回机器唯一标识符
// 计算MD5( 主机名 + 进程ID + 随机数 )
func GetMachineID() (string, error) {
	// 主机名
	host, err := os.Hostname()
	if nil != err {
		return "", err
	}
	// 进程ID
	pidstr := strconv.FormatInt(int64(os.Getpid()), 10)
	// 随机数
	uintByte := make([]byte, 4)
	binary.BigEndian.PutUint32(uintByte, uint32(rand.Int31()))
	randhex := hex.EncodeToString(uintByte)
	// 计算MD5
	machineid := GetMD5(strings.Join([]string{host, pidstr, randhex}, ","))
	return machineid, nil
}

// GetRandom 生成随机字符, 这个函数存在重复的概率, 需要唯一序列请使用GetUUID函数
func GetRandom(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

var regexpStrs = []string{"^", "$", ".", "*", "+", "?", "|", "/", "(", ")", "[", "]", "{", "}", "=", "!", ":", "-", ","}

// ReplaceRegexpSymbol 替换正则关键字
func ReplaceRegexpSymbol(str string) string {
	for i := 0; i < len(regexpStrs); i++ {
		str = strings.ReplaceAll(str, regexpStrs[i], "\\"+regexpStrs[i])
	}
	return str
}

// GetPlaceholder 生成占位符 '?', ',', 2 => '?,?'
func GetPlaceholder(placeholder, joinstr string, len int) (res string) {
	for i := 0; i < len; i++ {
		if i == 0 {
			res = placeholder
		} else {
			res += joinstr + placeholder
		}
	}
	return res
}

// ToInterface string类型转interface类型
func ToInterface(input ...string) []interface{} {
	output := make([]interface{}, len(input))
	if leni := len(input); leni > 0 {
		for i := 0; i < leni; i++ {
			output[i] = input[i]
		}
	}
	return output
}

// RemoveEmpty 去除空值
func RemoveEmpty(input ...string) []string {
	output := make([]string, 0)
	if len(input) > 0 {
		for _, str := range input {
			if len(str) > 0 {
				output = append(output, str)
			}
		}
	}
	return output
}

// RemoveDuplicatesAndEmpty 去重并去空
func RemoveDuplicatesAndEmpty(input ...string) []string {
	result := make([]string, 0)
	if len(input) > 0 {
		uniqueMap := make(map[string]interface{})
		for _, val := range input {
			if len(val) == 0 {
				continue
			}
			if _, ok := uniqueMap[val]; !ok {
				uniqueMap[val] = nil
				result = append(result, val)
			}
		}
	}
	return result
}

// Contain 是否在数组中包含
func Contain[T comparable](array []T, test T) bool {
	if len(array) == 0 {
		return false
	}
	for i := 0; i < len(array); i++ {
		if array[i] == test {
			return true
		}
	}
	return false
}

// ArrayMap 类似java stream.map
func ArrayMap[T any, R any](array []T, m func(row T) R) (res []R) {
	if len(array) == 0 {
		return
	}
	res = make([]R, len(array))
	for i := 0; i < len(array); i++ {
		res[i] = m(array[i])
	}
	return
}

// List2Array list对象转数组
func List2Array(list *list.List) []interface{} {
	if list.Len() == 0 {
		return nil
	}
	var arr []interface{}
	for e := list.Front(); e != nil; e = e.Next() {
		arr = append(arr, e.Value)
	}
	return arr
}
