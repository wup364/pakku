// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 配置工具-JSON文件实现
// 依赖包: utypes.Object fileutil

package jsonconfig

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/fileutil"
	"github.com/wup364/pakku/pkg/strutil"
	"github.com/wup364/pakku/pkg/utypes"
)

func init() {
	// 注册实例实现
	ipakku.PakkuConf.RegisterPakkuModuleImplement(new(Config), "IConfig", "json")
}

// Config json配置器
type Config struct {
	jsonObject map[string]any
	configPath string
	l          *sync.RWMutex
}

// Init 初始化解析器
func (config *Config) Init(appName string) error {
	path, err := filepath.Abs(".conf/" + appName + ".json")
	if nil != err {
		return err
	}
	return config.InitConfig(path)
}

// InitConfig 初始化解析器
func (config *Config) InitConfig(configPath string) error {
	if len(configPath) == 0 {
		return errors.New("config file path is empty")
	}
	// 创建父级目录
	parent := strutil.GetPathParent(configPath)
	if !fileutil.IsExist(parent) {
		err := fileutil.MkdirAll(parent)
		if nil != err {
			return err
		}
	}
	config.configPath = configPath
	// 文件不存在则创建
	if !fileutil.IsFile(config.configPath) {
		err := config.writeFileAsJSON(config.configPath, make(map[string]any))
		if nil != err {
			return err
		}
	}

	config.l = new(sync.RWMutex)
	config.l.Lock()
	defer config.l.Unlock()
	// Json to map
	config.jsonObject = make(map[string]any)
	return config.readFileAsJSON(config.configPath, &config.jsonObject)
}

// GetConfig 读取key的value信息
// 返回ConfigBody对象, 里面的值可能是string或者map
func (config *Config) GetConfig(key string) (res utypes.Object) {
	config.l.RLock()
	defer config.l.RUnlock()
	if len(key) == 0 || config.jsonObject == nil || len(config.jsonObject) == 0 {
		return
	}
	keys := strings.Split(key, ".")
	if keys == nil {
		return
	}
	var temp any
	keyLength := len(keys)
	for i := 0; i < keyLength; i++ {
		// last key
		if i == keyLength-1 {
			if i == 0 {
				if tp, ok := config.jsonObject[keys[i]]; ok {
					res = utypes.NewObject(tp)
				}
			} else if temp != nil {
				if tp, ok := temp.(map[string]any)[keys[i]]; ok {
					res = utypes.NewObject(tp)
				}
			}
			return
		}

		//
		var temp2 any
		if temp == nil { // first
			if tp, ok := config.jsonObject[keys[i]]; ok {
				temp2 = tp
			}
		} else { //
			if tp, ok := temp.(map[string]any)[keys[i]]; ok {
				temp2 = tp
			}
		}

		// find
		if temp2 != nil {
			temp = temp2
		} else {
			return
		}
	}
	return
}

// SetConfig 保存配置
func (config *Config) SetConfig(key string, value any) error {
	if len(key) == 0 || nil == value {
		return errors.New("key or value is empty")
	}
	config.l.Lock()
	defer config.l.Unlock()
	keys := strings.Split(key, ".")
	keyLength := len(keys)
	var temp any
	for i := 0; i < keyLength; i++ {
		// last key
		if i == keyLength-1 {
			if i == 0 {
				config.jsonObject[keys[i]] = value
			} else if temp != nil {
				temp.(map[string]any)[keys[i]] = value
			}
			err := config.writeFileAsJSON(config.configPath, config.jsonObject)
			return err
		}

		//
		var temp1 any
		if temp == nil { // first
			if tp, ok := config.jsonObject[keys[i]]; ok {
				temp1 = tp
			} else {
				temp1 = make(map[string]any)
				config.jsonObject[keys[i]] = temp1
			}
		} else { //
			if tp, ok := temp.(map[string]any)[keys[i]]; ok {
				temp1 = tp
			} else {
				temp1 = make(map[string]any)
				temp.(map[string]any)[keys[i]] = temp1
			}
		}

		// find
		if temp1 != nil {
			temp = temp1
		}
	}
	return nil
}

// readFileAsJSON 读取Json文件
func (config *Config) readFileAsJSON(path string, v any) error {
	if len(path) == 0 {
		return fileutil.PathNotExist("ReadFileAsJSON", path)
	}
	fp, err := os.OpenFile(path, os.O_RDONLY, 0)
	defer func() {
		if nil != fp {
			fp.Close()
		}
	}()

	if err == nil {
		decoder := json.NewDecoder(fp)
		decoder.UseNumber()
		return decoder.Decode(v)
	}
	return err
}

// writeFileAsJSON 写入Json文件
func (config *Config) writeFileAsJSON(path string, v any) error {
	if len(path) == 0 {
		return fileutil.PathNotExist("WriteFileAsJSON", path)
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
