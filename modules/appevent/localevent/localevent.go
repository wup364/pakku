package localevent

import (
	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/utypes"
)

func init() {
	ipakku.Override.RegisterInterfaceImpl(NewAppLocalEvent(), "IEvent", "local")
}

// NewAppLocalEvent NewAppLocalEvent
func NewAppLocalEvent() *AppLocalEvent {
	return &AppLocalEvent{smap: utypes.NewSafeMap()}
}

// AppLocalEvent 本机事件, 同步操作, 有结果返回
type AppLocalEvent struct {
	smap *utypes.SafeMap
}

// PublishSyncEvent PublishSyncEvent
func (ev *AppLocalEvent) PublishSyncEvent(group string, name string, val interface{}) error {
	var eventFunc ipakku.EventHandle
	if fun, ok := ev.smap.Get(group + name); ok {
		eventFunc = fun.(ipakku.EventHandle)
	} else {
		return ipakku.ErrSyncEventUnregistered
	}
	res := make(chan error)
	defer close(res)
	go (func(eventFunc ipakku.EventHandle, val interface{}, res chan error) {
		res <- eventFunc(val)
	})(eventFunc, val, res)
	return <-res
}

// ConsumerSyncEvent ConsumerSyncEvent
func (ev *AppLocalEvent) ConsumerSyncEvent(group string, name string, fun ipakku.EventHandle) error {
	if val, ok := ev.smap.Get(group + name); ok && nil != val {
		return ipakku.ErrSyncEventRegistered
	}
	ev.smap.Put(group+name, fun)
	return nil
}

// Init Init
func (ev *AppLocalEvent) Init(conf ipakku.AppConfig) error {
	return nil
}

// PublishEvent PublishEvent
func (ev *AppLocalEvent) PublishEvent(group string, name string, val interface{}) error {
	return ipakku.ErrEventMethodUnsupported
}

// ConsumerEvent ConsumerEvent
func (ev *AppLocalEvent) ConsumerEvent(group string, name string, fun ipakku.EventHandle) error {
	return ipakku.ErrEventMethodUnsupported
}
