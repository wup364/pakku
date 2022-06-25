// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
	fmt.Printf("改变后的值 %v\r\n", obj)
	obj.Upper.SayHello("...")
}
func TestSetStructFieldValueUnSafe(t *testing.T) {

	obj := &TestStruct{}
	val := &TestStructA{Upper: "泥马"}
	fmt.Printf("改变前的值 %v\r\n", obj)
	SetStructFieldValueUnSafe(obj, "lower", val)
	fmt.Printf("改变后的值 %v\r\n", obj)
	obj.lower.SayHello("...")
}
func TestSetInterfaceValueUnSafe(t *testing.T) {
	var obj SayHello
	val := TestStructA{Upper: "泥马"}
	fmt.Printf("改变前的值 %v\r\n", obj)
	if err := SetInterfaceValueUnSafe(&obj, &val); nil != err {
		panic(err)
	}
	fmt.Printf("改变后的值 %v\r\n", obj)
	obj.SayHello("...")
}
