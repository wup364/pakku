package appconfig

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/utypes"

	// 通过 init 函数注册
	_ "github.com/wup364/pakku/modules/appconfig/jsonconfig"
)

// AppConfig 配置模块
type AppConfig struct {
	configname string
	config     ipakku.IConfig
}

// AsModule 作为一个模块加载
func (conf *AppConfig) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "AppConfig",
		Version:     1.0,
		Description: "配置模块",
		OnReady: func(mctx ipakku.Loader) {
			// 获取配置的适配器, 默认json
			if err := ipakku.Override.AutowireInterfaceImpl(mctx, &conf.config, "json"); nil != err {
				logs.Panicln(err)
			}
			conf.configname = mctx.GetParam(ipakku.PARAMKEY_APPNAME).ToString("app")
		},
		OnSetup:  func() {},
		OnUpdate: func(cv float64) {},
		OnInit: func() {
			// 初始化配置
			conf.config.Init(conf.configname)
		},
	}
}

// GetConfig 读取key的value信息, 返回 Object 对象, 里面的值可能是string或者map
func (conf *AppConfig) GetConfig(key string) (res utypes.Object) {
	return conf.config.GetConfig(key)
}

// SetConfig 保存配置, key value 都为stirng
func (conf *AppConfig) SetConfig(key string, value string) error {
	return conf.config.SetConfig(key, value)
}
