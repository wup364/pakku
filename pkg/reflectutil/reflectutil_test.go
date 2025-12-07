// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// 反射工具

package reflectutil

import (
	"fmt"
	"testing"
)

type SayHello interface {
	SayHello(text string)
}

type TestStructA struct {
	Upper string
}

func (ta *TestStructA) SayHello(text string) {
	fmt.Println("艹" + ta.Upper)
}

type TestStruct struct {
	lower SayHello
	Upper SayHello
}

func TestSetStructFieldValue(t *testing.T) {
	obj := &TestStruct{}
	val := &TestStructA{}
	SetStructFieldValue(obj, "Upper", val)
	fmt.Printf("改变后的值 %v", obj)
	obj.Upper.SayHello("...")
}
func TestSetStructFieldValueUnSafe(t *testing.T) {

	obj := &TestStruct{}
	val := &TestStructA{Upper: "泥马"}
	fmt.Printf("改变前的值 %v", obj)
	SetStructFieldValueUnSafe(obj, "lower", val)
	fmt.Printf("改变后的值 %v", obj)
	obj.lower.SayHello("...")
}
func TestSetInterfaceValueUnSafe(t *testing.T) {
	var obj SayHello
	val := TestStructA{Upper: "泥马"}
	fmt.Printf("改变前的值 %v", obj)
	if err := SetInterfaceValueUnSafe(&obj, &val); nil != err {
		panic(err)
	}
	fmt.Printf("改变后的值 %v", obj)
	obj.SayHello("...")
}
