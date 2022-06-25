// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package serviceutil

import (
	"net/http"
	"testing"

	"github.com/wup364/pakku/utils/logs"
)

func TestHttpService(t *testing.T) {
	svr := NewHTTPService()
	svr.SetDebug(true)
	svr.Get("/hello/:[0-9]*$", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello : " + r.Method + ": " + r.URL.String()))
	})
	if err := svr.StartHTTP(StartHTTPConf{
		ListenAddr: "0.0.0.0:8080",
	}); nil != err {
		logs.Panicln(err)
	}
}
