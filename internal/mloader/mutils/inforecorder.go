// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 模块信息记录器-本地JSON文件实现
// 依赖包: utypes.Object fileutil

package mutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/fileutil"
	"github.com/wup364/pakku/pkg/strutil"
	"github.com/wup364/pakku/pkg/utypes"
)

func init() {
	ipakku.PakkuConf.SetModuleInfoRecorderImplement(new(InfoRecorder))
}

// InfoRecorder json配置器
type InfoRecorder struct {
	jsonObject map[string]any
	configPath string
	l          *sync.RWMutex
}

// Init 初始化解析器
func (config *InfoRecorder) Init(appName string) error {
	return config.InitConfig(fmt.Sprintf(".conf/%s_modules.json", appName))
}

// InitConfig 初始化解析器
func (config *InfoRecorder) InitConfig(configPath string) error {
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

// GetValue 读取配置
func (config *InfoRecorder) GetValue(key string) string {
	return config.GetValueByKey(key).ToString("")
}

// SetValue 写入配置
func (config *InfoRecorder) SetValue(key string, value string) error {
	return config.SetValueByKey(key, value)
}

// GetValueByKey 读取key的value信息
// 返回ModuleInfoRecorderBody对象, 里面的值可能是string或者map
func (config *InfoRecorder) GetValueByKey(key string) (res utypes.Object) {
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

// SetValueByKey 保存配置, key value 都为stirng
func (config *InfoRecorder) SetValueByKey(key string, value string) error {
	if len(key) == 0 || len(value) == 0 {
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
		var temp2 any
		if temp == nil { // first
			if tp, ok := config.jsonObject[keys[i]]; ok {
				temp2 = tp
			} else {
				temp2 = make(map[string]any)
				config.jsonObject[keys[i]] = temp2
			}
		} else { //
			if tp, ok := temp.(map[string]any)[keys[i]]; ok {
				temp2 = tp
			} else {
				temp2 = make(map[string]any)
				temp.(map[string]any)[keys[i]] = temp2
			}
		}

		// find
		if temp2 != nil {
			temp = temp2
		}
	}
	return nil
}

// readFileAsJSON 读取Json文件
func (config *InfoRecorder) readFileAsJSON(path string, v any) error {
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

// writeFileAsJSON 写入Json文件
func (config *InfoRecorder) writeFileAsJSON(path string, v any) error {
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
