package localevent

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/utypes"
)

func init() {
	ipakku.PakkuConf.RegisterPakkuModuleImplement(NewAppLocalEvent(), "IEvent", "local")
}

// NewAppLocalEvent NewAppLocalEvent
func NewAppLocalEvent() *AppLocalEvent {
	return &AppLocalEvent{smap: utypes.NewSafeMap[string, []ipakku.EventHandle]()}
}

// AppLocalEvent 本机事件, 同步操作, 有结果返回
type AppLocalEvent struct {
	smap *utypes.SafeMap[string, []ipakku.EventHandle]
}

// PublishSyncEvent PublishSyncEvent
func (ev *AppLocalEvent) PublishSyncEvent(group string, name string, val any) (err error) {
	var eventFuncs []ipakku.EventHandle
	if fun, ok := ev.smap.Get(group + name); ok {
		eventFuncs = fun
	} else {
		logs.Errorf("event unregistered: group=%s, name=%s \r\n", group, name)
		return ipakku.ErrSyncEventUnregistered
	}

	//
	for i := 0; i < len(eventFuncs); i++ {
		if err = eventFuncs[i](val); nil != err {
			return
		}
	}
	return
}

// ConsumerSyncEvent ConsumerSyncEvent
func (ev *AppLocalEvent) ConsumerSyncEvent(group string, name string, fun ipakku.EventHandle) (err error) {
	if nil == fun {
		return
	}
	var eventFuncs []ipakku.EventHandle
	if val, ok := ev.smap.Get(group + name); ok && nil != val {
		eventFuncs = val
	} else {
		eventFuncs = make([]ipakku.EventHandle, 0)
	}
	ev.smap.Put(group+name, append(eventFuncs, fun))
	return
}

// Init Init
func (ev *AppLocalEvent) Init(conf ipakku.AppConfig) error {
	return nil
}

// PublishEvent PublishEvent
func (ev *AppLocalEvent) PublishEvent(group string, name string, val any) error {
	return ipakku.ErrEventMethodUnsupported
}

// ConsumerEvent ConsumerEvent
func (ev *AppLocalEvent) ConsumerEvent(group string, name string, fun ipakku.EventHandle) error {
	return ipakku.ErrEventMethodUnsupported
}
