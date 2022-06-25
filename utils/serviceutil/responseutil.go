package serviceutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/wup364/pakku/utils/fileutil"
	"github.com/wup364/pakku/utils/strutil"
)

// HTTPResponse 接口返回格式约束
type HTTPResponse struct {
	Code int         `json:"code"`
	Flag string      `json:"flag"`
	Data interface{} `json:"data"`
}

// SendSuccess 返回成功结果
func SendSuccess(w http.ResponseWriter, msg interface{}) {
	SendSuccessAndStatus(w, http.StatusOK, msg)
}

// SendBadRequest 返回400错误
func SendBadRequest(w http.ResponseWriter, msg interface{}) {
	SendErrorAndStatus(w, http.StatusBadRequest, msg)
}

// SendServerError 返回500错误
func SendServerError(w http.ResponseWriter, msg interface{}) {
	SendErrorAndStatus(w, http.StatusInternalServerError, msg)
}

// SendSuccessAndStatus 返回成功结果
func SendSuccessAndStatus(w http.ResponseWriter, statusCode int, msg interface{}) {
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(BuildHttpResponse(statusCode, "T", msg))
}

// SendErrorAndStatus 返回失败结果
func SendErrorAndStatus(w http.ResponseWriter, statusCode int, msg interface{}) {
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(BuildHttpResponse(statusCode, "F", msg))
}

// BuildHttpResponse 构建返回json
func BuildHttpResponse(code int, flag string, str interface{}) []byte {
	bt, err := json.Marshal(HTTPResponse{Code: code, Flag: flag, Data: str})
	if nil != err {
		return []byte(err.Error())
	}
	return bt
}

// 将请求体解析为对象
func ParseHTTPRequest(r *http.Request, obj interface{}) error {
	return json.Unmarshal([]byte(strutil.ReadAsString(r.Body)), obj)
}

// Parse2HTTPResponse json转对象
func Parse2HTTPResponse(str string) *HTTPResponse {
	res := &HTTPResponse{}
	if err := json.Unmarshal([]byte(str), res); nil != err {
		return nil
	}
	return res
}

// WirteFile 发送文件流, 支持分段
func WirteFile(w http.ResponseWriter, r *http.Request, path string) {
	// 校验
	if !fileutil.IsFile(path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if maxSize, err := fileutil.GetFileSize(path); err != nil {
		SendServerError(w, err.Error())
	} else {
		var sa *os.File
		if sa, err = fileutil.OpenFile(path); nil != err {
			SendServerError(w, err.Error())
		} else {
			defer sa.Close()
			start, end, hasRange := GetRequestRange(r, maxSize)
			RangeWrite(w, sa, start, end, maxSize, hasRange)
		}
	}
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

// RangeWrite 范围写入http, 如文件分段传输
func RangeWrite(w http.ResponseWriter, sa io.ReadSeeker, start, end, maxSize int64, hasRange bool) {
	if _, err := sa.Seek(start, io.SeekStart); nil != err {
		SendServerError(w, err.Error())
	} else {
		ctLength := end - start
		if ctLength == 0 || ctLength < 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Length", strconv.Itoa(int(ctLength)))
			if hasRange {
				w.Header().Set("Content-Range", "bytes "+strconv.Itoa(int(start))+"-"+strconv.Itoa(int(end-1))+"/"+strconv.Itoa(int(maxSize)))
				w.WriteHeader(http.StatusPartialContent)
			}
			if _, err := io.Copy(w, io.LimitReader(sa, ctLength)); nil != err && err != io.EOF {
				SendServerError(w, err.Error())
			}
		}
	}
}
