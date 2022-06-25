package ipakku

import "errors"

// EventHandle 异步事件回调
type EventHandle func(v interface{}) error

// ErrSyncEventUnregistered 事件未注册
var ErrSyncEventUnregistered = errors.New("sync event unregistered")

// ErrSyncEventRegistered 事件重复注册
var ErrSyncEventRegistered = errors.New("sync event is registered")

// ErrEventMethodUnsupported 没有实现
var ErrEventMethodUnsupported = errors.New("event method unsupported")

// AppEvent 事件模型
type AppEvent interface {
	// 异步事件, 默认未实现
	PublishEvent(group string, name string, val interface{}) error
	ConsumerEvent(group string, name string, fun EventHandle) error
}

// AppSyncEvent 本机事件[不开放自定义实现], 同步操作 只能注册一次
type AppSyncEvent interface {
	PublishSyncEvent(group string, name string, val interface{}) error
	ConsumerSyncEvent(group string, name string, fun EventHandle) error
}
