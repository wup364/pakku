// SPDX-License-Identifier: MIT
// Copyright (C) 2025 WuPeng <wup364@outlook.com>.

package appstaticpage

import "github.com/wup364/pakku/pkg/fileutil"

const (
	// DefaultStaticPageConfigPath 默认的WebPage配置文件路径
	DefaultStaticPageConfigPath = ".conf/pakku-static-pages.json"
)

// GetStaticPageConfig 获取 WebPage 的配置
// 如果配置文件不存在，则返回 nil
func GetStaticPageConfig() (*StaticPageConfig, error) {
	if !fileutil.IsFile(DefaultStaticPageConfigPath) {
		return nil, nil
	}
	var result StaticPageConfig
	if err := fileutil.ReadFileAsJSON(DefaultStaticPageConfigPath, &result); nil != err {
		return nil, err
	} else {
		return &result, nil
	}
}

// StaticPageConfig 表示静态页面配置
type StaticPageConfig struct {
	EnableCORS      bool            `json:"enableCORS"`
	Redirect        Redirect        `json:"redirect"`
	StaticFiles     []StaticFile    `json:"staticFiles"`
	AllowedMethods  []string        `json:"allowedMethods"`
	AllowedHeaders  []string        `json:"allowedHeaders"`
	AllowedOrigins  []string        `json:"allowedOrigins"`
	StaticDirectory StaticDirectory `json:"staticDirectory"`
}

// Redirect 表示重定向配置
type Redirect struct {
	Path   string `json:"path"`
	Target string `json:"target"`
	Status int    `json:"status"`
}

// StaticFile 表示静态文件配置
type StaticFile struct {
	Path     string `json:"path"`
	FilePath string `json:"filePath"`
}

// StaticDirectory 表示静态目录配置
type StaticDirectory struct {
	Path      string `json:"path"`
	Directory string `json:"directory"`
}
