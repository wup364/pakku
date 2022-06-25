package appevent

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"

	"github.com/wup364/pakku/modules/appevent/localevent"
)

// AppEvent 事件模块
type AppEvent struct {
	event  ipakku.AppEvent
	sysevt ipakku.AppSyncEvent
	conf   ipakku.AppConfig `@autowired:"AppConfig"`
}

// AsModule 作为一个模块加载
func (ev *AppEvent) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "AppEvent",
		Version:     1.0,
		Description: "事件模块",
		OnReady: func(mctx ipakku.Loader) {
			var driver ipakku.IEvent
			if err := ipakku.Override.AutowireInterfaceImpl(mctx, &driver, "local"); nil != err {
				logs.Panicln(err)
			} else if err := driver.Init(ev.conf); nil != err {
				logs.Panicln(err)
			}
			ev.event = driver
			ev.sysevt = localevent.NewAppLocalEvent()
		},
		OnSetup:  func() {},
		OnUpdate: func(cv float64) {},
		OnInit: func() {
		},
	}
}

// PublishSyncEvent PublishSyncEvent
func (ev *AppEvent) PublishSyncEvent(group string, name string, val interface{}) error {
	return ev.sysevt.PublishSyncEvent(group, name, val)
}

// ConsumerSyncEvent ConsumerSyncEvent
func (ev *AppEvent) ConsumerSyncEvent(group string, name string, fun ipakku.EventHandle) error {
	return ev.sysevt.ConsumerSyncEvent(group, name, fun)
}

// PublishEvent PublishEvent
func (ev *AppEvent) PublishEvent(name string, val string, obj interface{}) error {
	return ev.event.PublishEvent(name, val, obj)
}

// ConsumerEvent ConsumerEvent
func (ev *AppEvent) ConsumerEvent(group string, name string, fun ipakku.EventHandle) error {
	return ev.event.ConsumerEvent(group, name, fun)
}
