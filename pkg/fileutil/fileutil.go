// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 文件工具

package fileutil

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/wup364/pakku/pkg/strutil"
)

// WalkCallback (源路径, 目标路径, 错误信息) 错误. 子文件操作回调
// 一般用于子文件操作时间较长的操作, 如: MoveFilesAcrossDisk, MoveFilesByCopying, CopyFiles. 返回 非空error则终止下一步操作
type WalkCallback func(src, dst string, err error) error

// GetFileInfo 获取文件信息对象
func GetFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// OpenFile 获取文件信息对象
func OpenFile(path string) (*os.File, error) {
	return os.Open(path)
}

// GetWriter 获取只写文件对象
// O_RDONLY: 只读模式(read-only)
// O_WRONLY: 只写模式(write-only)
// O_RDWR: 读写模式(read-write)
// O_APPEND: 追加模式(append)
// O_CREATE: 文件不存在就创建(create a new file if none exists.)
// O_EXCL: 与 O_CREATE 一起用, 构成一个新建文件的功能, 它要求文件必须不存在
// O_SYNC: 同步方式打开，即不使用缓存，直接写入硬盘
// O_TRUNC: 打开并清空文件
func GetWriter(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

// IsExist 文件/夹是否存在
func IsExist(path string) bool {
	_, err := GetFileInfo(path)
	return err == nil
}

// IsDir 是否是文件夹
func IsDir(path string) bool {
	stat, err := GetFileInfo(path)
	if err == nil {
		return stat.IsDir()
	}
	return false
}

// MkdirAll 创建文件夹-多级
func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// Mkdir 创建文件夹
func Mkdir(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

// IsFile 是否是文件
func IsFile(path string) bool {
	stat, err := GetFileInfo(path)
	if err == nil {
		return !stat.IsDir()
	}
	return false
}

// GetDirList 获取一级子目录名字(包含文件|文件夹,无序)
func GetDirList(path string) ([]string, error) {
	f, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	list, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetDirListWithFolderFirst 获取一级子目录名字(包含文件|文件夹,无序,文件夹在前)
func GetDirListWithFolderFirst(path string) ([]string, error) {
	if files, dirs, err := GetDirAndFileList(path); nil != err {
		return nil, err
	} else {
		return append(dirs, files...), nil
	}
}

// GetDirAndFileList 获取一级子目录名字(包含文件|文件夹,无序, 文件和文件夹分开)
func GetDirAndFileList(path string) (files, dirs []string, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()

	// 读取目录内容
	var infos []os.FileInfo
	if infos, err = f.Readdir(-1); err != nil {
		return
	}

	// 分离文件和目录
	for _, info := range infos {
		if info.IsDir() {
			dirs = append(dirs, info.Name())
		} else {
			files = append(files, info.Name())
		}
	}
	return
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	f, err := OpenFile(path)
	if err != nil {
		return 0, err
	}
	fInfo, errStat := f.Stat()
	f.Close()
	if errStat != nil {
		return 0, errStat
	}
	return fInfo.Size(), nil
}

// GetCreateTime 获取创建时间
func GetCreateTime(path string) (time.Time, error) {
	return GetModifyTime(path)
}

// GetModifyTime 获取修改时间
func GetModifyTime(path string) (time.Time, error) {
	f, err := OpenFile(path)
	if err != nil {
		return time.Time{}, err
	}
	fInfo, errStat := f.Stat()
	f.Close()
	if errStat != nil {
		return time.Time{}, errStat
	}
	return fInfo.ModTime(), nil
}

// RemoveFile 删除文件
func RemoveFile(file string) error {
	if !IsExist(file) {
		return PathNotExist("RemoveFile", file)
	}
	return os.Remove(file)
}

// RemoveAll 删除文件
func RemoveAll(file string) error {
	if !IsExist(file) {
		return PathNotExist("RemoveAll", file)
	}
	return os.RemoveAll(file)
}

// Rename 重命名
func Rename(old, newName string) error {
	path := path.Clean(old)
	if strings.Contains(path, "\\") {
		return os.Rename(old, old[:strings.LastIndex(old, "\\")+1]+newName)
	}
	return os.Rename(old, old[:strings.LastIndex(old, "/")+1]+newName)
}

// MoveFilesAcrossDisk 移动文件|文件夹,可跨分区移动(源路径, 目标路径, 重复覆盖, 重复忽略, 每个路径操作结果回调(可以nil))
func MoveFilesAcrossDisk(src, dst string, replace, jump bool, walk WalkCallback) error {
	src = path.Clean(src)
	dst = path.Clean(dst)
	if src == dst && len(src) > 0 {
		return nil
	}
	if !IsExist(src) {
		return PathNotExist("MoveFilesAcrossDisk", src)
	}
	// 尝试本分区移动
	err := MoveFiles(src, dst, replace, jump)
	// 尝试跨分区移动
	if IsAcrossDiskError(err) {
		return MoveFilesByCopying(src, dst, replace, jump, walk)
	}
	return err
}

// MoveFilesByCopying 移动文件|夹(拷贝+删除) (源路径, 目标路径, 重复覆盖, 重复忽略, 每个路径操作结果回调(可以nil))
func MoveFilesByCopying(src, dst string, replace, jump bool, walk WalkCallback) error {
	src = path.Clean(src)
	dst = path.Clean(dst)
	doWalk := func(src string, dst string, err error) error {
		if nil == walk {
			return err
		}
		return walk(src, dst, err)
	}
	if src == dst && len(src) > 0 {
		return nil
	}
	if IsFile(src) {
		err := CopyFile(src, dst, replace, jump)
		if nil == err {
			return os.Remove(src)
		}
		return doWalk(src, dst, err)
	}
	// 复制文件夹
	err := CopyFiles(src, dst, replace, jump, func(src string, dst string, err error) error {
		return doWalk(src, dst, err)
	})
	// 最后的清理
	if nil == err {
		err = os.RemoveAll(src)
	}
	return err
}

// CopyFiles 复制文件夹 (源路径, 目标路径, 重复覆盖, 重复忽略, 每个路径操作结果回调(可以nil))
// walk 返回err则停止下一个拷贝
func CopyFiles(src, dst string, replace, jump bool, walk WalkCallback) error {
	if src == dst && len(src) > 0 {
		return nil
	}
	if !IsExist(src) || !IsDir(src) {
		return PathNotExist("CopyFiles", src)
	}
	if IsFile(dst) {
		return PathExist("CopyFiles", dst)
	}
	doWalk := func(src string, dst string, err error) error {
		if nil == walk {
			return err
		}
		return walk(src, dst, err)
	}
	srcLen := len(src)
	return filepath.Walk(src, func(s string, f os.FileInfo, err error) error {
		d := dst + s[srcLen:]
		if err == nil {
			if f.IsDir() {
				if IsFile(d) {
					if replace {
						err = os.Remove(d)
					} else if !jump {
						return PathExist("CopyFiles", d)
					}
				}
				if nil != err {
					return err
				} else if !IsExist(d) {
					err = os.Mkdir(d, os.ModePerm)
				}
				return err
			}
			return doWalk(s, d, CopyFile(s, d, replace, jump))
		}
		return err
	})
}

// MoveFiles 移动文件|夹, 如果存在的话就列表后逐一移动 (源路径, 目标路径, 重复覆盖, 重复忽略)
func MoveFiles(src, dst string, replace, jump bool) error {
	if src == dst && len(src) > 0 {
		return nil
	}
	if !IsExist(src) {
		return PathNotExist("MoveFiles", src)
	}
	srcIsfile := IsFile(src)
	if IsExist(dst) {
		dstIsfile := IsFile(dst)
		var err error
		// src&dst都是文件夹, 则考虑合并规则, 处理里面的文件
		if !srcIsfile && !dstIsfile {
			list, _ := GetDirList(src)
			for i := 0; i < len(list); i++ {
				err = MoveFiles(src+"/"+list[i], dst+"/"+list[i], replace, jump)
				if err != nil {
					return err // 返回终止信号
				}
			}
			if IsExist(src) {
				err = os.RemoveAll(src)
			}
			return err
		}
		// 非合并模式只考虑删除目标位置
		if jump {
			return err
		} else if replace {
			// 如果覆盖, 则先删除目标位置
			if dstIsfile {
				err = os.Remove(dst)
			} else {
				err = os.RemoveAll(dst)
			}
		} else {
			err = PathExist("MoveFiles", dst)
		}
		if nil != err {
			return err
		}
		return os.Rename(src, dst)
	}
	{
		dstParent := strutil.GetPathParent(dst)
		if !IsExist(dstParent) {
			if err := MkdirAll(dstParent); nil != err {
				return err
			}
		}
	}
	// 如果目标文件/夹不存在就直接移动
	return os.Rename(src, dst)
}

// CopyFile 复制文件(源路径, 目标路径, 重复覆盖, 重复忽略)
func CopyFile(src, dst string, replace, jump bool) error {
	if src == dst && len(src) > 0 {
		return nil
	}
	if !IsFile(src) {
		return PathNotExist("CopyFile", src)
	}
	if IsExist(dst) {
		var err error
		if replace {
			if IsFile(src) {
				err = os.Remove(dst)
			} else {
				err = os.RemoveAll(dst)
			}
		} else if !jump {
			err = PathExist("CopyFile", dst)
		}
		if err != nil {
			return err
		}
	}
	dstParent := strutil.GetPathParent(dst)
	if !IsExist(dstParent) {
		if err := MkdirAll(dstParent); err != nil {
			return err
		}
	}
	RSrc, err := OpenFile(src)
	defer func() {
		if nil != RSrc {
			RSrc.Close()
		}
	}()
	if err != nil {
		return err
	}
	WDst, err := GetWriter(dst)
	defer func() {
		if nil != WDst {
			WDst.Close()
		}
	}()
	if err != nil {
		return err
	}
	_, err = io.Copy(WDst, RSrc)
	return err
}

// ReadFileAsJSON 读取Json文件
func ReadFileAsJSON(path string, v any) error {
	if len(path) == 0 {
		return PathNotExist("ReadFileAsJSON", path)
	}
	fp, err := os.OpenFile(path, os.O_RDONLY, 0)
	defer func() {
		if nil != fp {
			fp.Close()
		}
	}()

	if err == nil {
		st, stErr := fp.Stat()
		if stErr == nil {
			data := make([]byte, st.Size())
			_, err = fp.Read(data)
			if err == nil {
				return json.Unmarshal(data, v)
			}
		} else {
			err = stErr
		}
	}
	return err
}

// WriteFileAsJSON 写入Json文件
func WriteFileAsJSON(path string, v any) error {
	if len(path) == 0 {
		return PathNotExist("WriteFileAsJSON", path)
	}
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	defer func() {
		if nil != fp {
			fp.Close()
		}
	}()

	if err == nil {
		data, err := json.Marshal(v)
		if err == nil {
			_, err := fp.Write(data)
			return err
		}
		return err
	}
	return err
}

// ReadFileAsText 读取Json文件
func ReadFileAsText(path string) (string, error) {
	if len(path) == 0 {
		return "", PathNotExist("ReadFileAsText", path)
	}
	fp, err := os.OpenFile(path, os.O_RDONLY, 0)
	defer func() {
		if nil != fp {
			fp.Close()
		}
	}()

	if err == nil {
		st, stErr := fp.Stat()
		if stErr == nil {
			data := make([]byte, st.Size())
			_, err = fp.Read(data)
			if err == nil {
				return string(data), nil
			}
		} else {
			err = stErr
		}
	}
	return "", err
}

// WriteTextFile 写入文本文件
func WriteTextFile(path, text string) error {
	if len(path) == 0 {
		return PathNotExist("WriteFile", path)
	}
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	defer func() {
		if nil != fp {
			fp.Close()
		}
	}()

	if err == nil {
		if err == nil {
			_, err := fp.Write([]byte(text))
			return err
		}
		return err
	}
	return err
}

// GetSHA256 获取sha256
func GetSHA256(r io.Reader) (string, error) {
	h := sha256.New()
	buf := make([]byte, 1<<20)
	for {
		n, err := io.ReadFull(r, buf)
		if err == nil || err == io.ErrUnexpectedEOF || err == io.EOF {
			if n > 0 {
				if _, err = h.Write(buf[0:n]); err != nil {
					return "", err
				}
			} else {
				break
			}
		} else {
			return "", err
		}
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// FormatFileSize 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// PathExist 路径已经存在的错误
func PathExist(op, path string) error {
	return &os.PathError{
		Op:   op,
		Path: path,
		Err:  os.ErrExist,
	}
}

// PathNotExist 路径不存在的错误
func PathNotExist(op, path string) error {
	return &os.PathError{
		Op:   op,
		Path: path,
		Err:  os.ErrNotExist,
	}
}

// IsExistError 是否是目标位置已经存在的错误
func IsExistError(err error) bool {
	if nil == err {
		return false
	}
	var cpErr error
	switch err := err.(type) {
	case *os.PathError:
		cpErr = err.Err
	}
	if cpErr == nil {
		return false
	}
	return cpErr.Error() == os.ErrExist.Error()
}

// IsNotExistError 是否是目标位置不存在的错误
func IsNotExistError(err error) bool {
	var cpErr error
	switch err := err.(type) {
	case *os.PathError:
		cpErr = err.Err
	}
	if cpErr == nil {
		return false
	}
	return cpErr.Error() == os.ErrNotExist.Error()
}

// IsAcrossDiskError 是否是跨磁盘错误
func IsAcrossDiskError(err error) bool {
	var cpErr error
	switch err := err.(type) {
	case *os.LinkError:
		cpErr = err.Err
	}
	if cpErr == nil {
		return false
	}
	return cpErr.Error() == "The system cannot move the file to a different disk drive."
}
