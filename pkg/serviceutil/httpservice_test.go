// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

package serviceutil

import (
	"net/http"
	"testing"

	"github.com/wup364/pakku/pkg/logs"
)

func TestHttpService(t *testing.T) {
	svr := NewHTTPService()
	svr.SetDebug(true)
	svr.Get("/hello/:*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello : " + r.Method + ": " + r.URL.String()))
	})
	if err := svr.StartHTTP(StartHTTPConf{
		ListenAddr: "0.0.0.0:8080",
	}); nil != err {
		logs.Panic(err)
	}
}
