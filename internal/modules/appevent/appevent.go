package appevent

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/pkg/logs"

	"github.com/wup364/pakku/internal/modules/appevent/localevent"
)

// AppEvent 事件模块
type AppEvent struct {
	event  ipakku.AppEvent
	sysevt ipakku.AppSyncEvent
	conf   ipakku.AppConfig `@autowired:""`
}

// AsModule 作为一个模块加载
func (ev *AppEvent) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Version:     1.0,
		Description: "AppEvent module",
		OnReady: func(app ipakku.Application) {
			var driver ipakku.IEvent
			if err := ipakku.PakkuConf.AutowirePakkuModuleImplement(app.Params(), &driver, "local"); nil != err {
				logs.Panic(err)
			} else if err := driver.Init(ev.conf); nil != err {
				logs.Panic(err)
			}
			ev.event = driver
			ev.sysevt = localevent.NewAppLocalEvent()
		},
	}
}

// PublishSyncEvent PublishSyncEvent
func (ev *AppEvent) PublishSyncEvent(group string, name string, val any) error {
	return ev.sysevt.PublishSyncEvent(group, name, val)
}

// ConsumerSyncEvent ConsumerSyncEvent
func (ev *AppEvent) ConsumerSyncEvent(group string, name string, fun ipakku.EventHandle) error {
	return ev.sysevt.ConsumerSyncEvent(group, name, fun)
}

// PublishEvent PublishEvent
func (ev *AppEvent) PublishEvent(name string, val string, obj any) error {
	return ev.event.PublishEvent(name, val, obj)
}

// ConsumerEvent ConsumerEvent
func (ev *AppEvent) ConsumerEvent(group string, name string, fun ipakku.EventHandle) error {
	return ev.event.ConsumerEvent(group, name, fun)
}
