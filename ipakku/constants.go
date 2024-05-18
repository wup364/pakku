package ipakku

//
const (

	// STAG_AUTOWIRED struct标签-自动注入标签
	STAG_AUTOWIRED = "@autowired"

	// STAG_AUTOCONFIG struct标签-自动配置标签
	STAG_AUTOCONFIG = "@autoConfig"

	// STAG_CONFIG_VALUE struct标签-自动配置-字段配置标签
	STAG_CONFIG_VALUE = "@value"
)

// ModuleID 模块ID
var ModuleID = moduleID{
	AppConfig:  "AppConfig",
	AppCache:   "AppCache",
	AppEvent:   "AppEvent",
	AppService: "AppService",
}

// moduleID 模块ID
type moduleID struct {
	AppConfig  string
	AppCache   string
	AppEvent   string
	AppService string
}
