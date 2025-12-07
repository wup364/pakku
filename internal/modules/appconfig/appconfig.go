package appconfig

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"
	"github.com/wup364/pakku/pkg/utypes"

	// 通过 init 函数注册
	"github.com/wup364/pakku/internal/modules/appconfig/confutils"
	_ "github.com/wup364/pakku/internal/modules/appconfig/jsonconfig"
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
		Version:     1.0,
		Description: "AppConfig module",
		OnReady: func(app ipakku.Application) {
			// 获取配置的适配器, 默认json
			if err := ipakku.PakkuConf.AutowirePakkuModuleImplement(app.Params(), &conf.config, "json"); nil != err {
				logs.Panic(err)
			}
			conf.configname = app.Params().GetParam(ipakku.PARAMS_KEY_APPNAME).ToString(ipakku.DEFT_VAL_APPNAME)
			conf.autoValue = confutils.NewAutoValueOfBeanUtil(conf.config)

			// 注册监听 - 自动完成配置类的配置
			app.Modules().OnModuleEvent("*", ipakku.ModuleEventOnReady, func(module any, app ipakku.Application) {
				if err := conf.ScanAndAutoConfig(module); nil != err {
					logs.Panic(err)
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
func (conf *AppConfig) SetConfig(key string, value any) error {
	return conf.config.SetConfig(key, value)
}

// ScanAndAutoConfig 扫描带有@autoconfig标签的字段, 并完成其配置
func (conf *AppConfig) ScanAndAutoConfig(ptr any) error {
	return conf.autoValue.ScanAndAutoConfig(ptr)
}

// ScanAndAutoValue 扫描带有@value标签的字段, 并完成其配置
func (conf *AppConfig) ScanAndAutoValue(configPrefix string, ptr any) error {
	return conf.autoValue.ScanAndAutoValue(configPrefix, ptr)
}
