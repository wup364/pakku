// Copyright (C) 2019 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// HTTP客户端工具

package httpclient

import (
	"errors"
	"net/http"
	"testing"

	"github.com/wup364/pakku/utils/strutil"
)

func TestGet(t *testing.T) {
	url := "http://127.0.0.1:8080/file/v1/list"
	header := map[string]string{
		"X-Ack":  "1d4116dd67902bc670c00704bb5a8581",
		"X-Sign": "4e44ad3632d8aca624f3022e7a0bc98b442f7de3aa04ffaf61a0fc30c4dc6260",
	}
	if resp, err := Get(url, map[string]string{"path": "/"}, header); nil == err {
		t.Logf("TestGet Result %s \r\n", strutil.ReadAsString(resp.Body))
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
