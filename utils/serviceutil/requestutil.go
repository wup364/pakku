package serviceutil

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/wup364/pakku/utils/strutil"
)

// 将请求体解析为对象
func ParseHTTPRequest(r *http.Request, obj interface{}) error {
	return json.Unmarshal([]byte(strutil.ReadAsString(r.Body)), obj)
}

// GetRequestRange 解析http分段头信息
func GetRequestRange(r *http.Request, maxSize int64) (start, end int64, hasRange bool) {
	var qRange string
	if qRange = r.Header.Get("Range"); len(qRange) == 0 {
		qRange = r.FormValue("Range")
	}
	if len(qRange) > 0 {
		hasRange = true
		temp := qRange[strings.Index(qRange, "=")+1:]
		if index := strings.Index(temp, "-"); index > -1 {
			var err error
			if start, err = strconv.ParseInt(temp[0:strings.Index(temp, "-")], 10, 64); nil != err || start < 0 {
				start = 0
			}
			if end, err = strconv.ParseInt(temp[strings.Index(temp, "-")+1:], 10, 64); nil != err || end == 0 {
				end = maxSize
			}
		}
	} else {
		end = maxSize
	}
	return start, end, hasRange
}
