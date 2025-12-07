// SPDX-License-Identifier: MIT
// Copyright (C) 2019 WuPeng <wup364@outlook.com>.

// HTTP客户端工具

package httpclient

import (
	"errors"
	"net/http"
	"testing"

	"github.com/wup364/pakku/pkg/strutil"
)

func TestGet(t *testing.T) {
	url := "http://127.0.0.1:8080/file/v1/list"
	header := map[string]string{
		"X-Ack":  "1d4116dd67902bc670c00704bb5a8581",
		"X-Sign": "4e44ad3632d8aca624f3022e7a0bc98b442f7de3aa04ffaf61a0fc30c4dc6260",
	}
	if resp, err := Get(url, map[string]string{"path": "/"}, header); nil == err {
		t.Logf("TestGet Result %s ", strutil.ReadAsString(resp.Body))
	} else {
		t.Error(err)
	}
}

func TestPostFile(t *testing.T) {
	url := "http://127.0.0.1:8080/filestream/v1/put/937e16dd6ce020f5897466cb3908d2ac"
	header := map[string]string{
		"FormName-File": "file",
	}
	if response, err := PostFile(url, "./httpclient.go", header); nil == err {
		if response.StatusCode == http.StatusOK {
			t.Log("TestPostFile ok")
		} else {
			t.Error(errors.New("[" + response.Status + "] " + strutil.ReadAsString(response.Body)))
		}
	} else {
		t.Error(err)
	}
}
