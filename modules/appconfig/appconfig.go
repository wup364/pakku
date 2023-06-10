package appconfig

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/utypes"

	// 通过 init 函数注册
	"github.com/wup364/pakku/modules/appconfig/confutils"
	_ "github.com/wup364/pakku/modules/appconfig/jsonconfig"
)

// AppConfig 配置模块
type AppConfig struct {
	configname string
	config     ipakku.IConfig
	autoValue  *confutils.AutoValueOfBeanUtil
}

// AsModule 作为一个模块加载
func (conf *AppConfig) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "AppConfig",
		Version:     1.0,
		Description: "AppConfig module",
		OnReady: func(mctx ipakku.Loader) {
			// 获取配置的适配器, 默认json
			if err := ipakku.Override.AutowireInterfaceImpl(mctx, &conf.config, "json"); nil != err {
				logs.Panicln(err)
			}
			conf.configname = mctx.GetParam(ipakku.PARAMKEY_APPNAME).ToString("app")
			conf.autoValue = confutils.NewAutoValueOfBeanUtil(conf.config)

			// 注册监听 - 自动完成配置类的配置
			mctx.OnModuleEvent("*", ipakku.ModuleEventOnReady, func(module interface{}, loader ipakku.Loader) {
				if err := conf.AutoValueOfBean(module); nil != err {
					logs.Panicln(err)
				}
			})
		},
		OnInit: func() {
			conf.config.Init(conf.configname)
		},
	}
}

// GetConfig 读取key的value信息, 返回 Object 对象, 里面的值可能是string或者map
func (conf *AppConfig) GetConfig(key string) (res utypes.Object) {
	return conf.config.GetConfig(key)
}

// SetConfig 设置值
func (conf *AppConfig) SetConfig(key string, value interface{}) error {
	return conf.config.SetConfig(key, value)
}

// AutoValueOfBean 根据bean描述自动完成字段值
func (conf *AppConfig) AutoValueOfBean(ptr interface{}) (err error) {
	return conf.autoValue.AutoValueOfBean(ptr)
}
