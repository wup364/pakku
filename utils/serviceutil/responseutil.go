package serviceutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/wup364/pakku/utils/fileutil"
	"github.com/wup364/pakku/utils/utypes"
)

// HTTPResponse 接口返回格式约束
type HTTPResponse struct {
	Code string      `json:"code"`
	Flag string      `json:"flag"`
	Data interface{} `json:"data"`
}

// BusinessError 业务异常
type BusinessError interface {
	error
	// GetErrorCode 获取业务异常代码
	GetErrorCode() string
}

// SendSuccess 返回成功结果 httpCode=200, code=OK
func SendSuccess(w http.ResponseWriter, msg interface{}) {
	SendSuccessResponse(w, http.StatusOK, "OK", msg)
}

// SendBusinessError 返回业务错误 httpCode=200,
// code值为默认为BUSINESS_ERROR, 若error实现了CustomError接口, 则使用CustomError.ErrorCode值
func SendBusinessError(w http.ResponseWriter, err error) {
	if cr, ok := err.(utypes.CustomError); ok && len(cr.ErrorCode()) > 0 {
		SendBusinessErrorAndCode(w, cr.ErrorCode(), err.Error())
	} else {
		SendBusinessErrorAndCode(w, "BUSINESS_ERROR", err.Error())
	}
}

// SendBadRequest 返回400错误, code=BAD_REQUEST
func SendBadRequest(w http.ResponseWriter, msg interface{}) {
	SendErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", msg)
}

// SendServerError 返回500错误, code=SERVER_ERROR
func SendServerError(w http.ResponseWriter, msg interface{}) {
	SendErrorResponse(w, http.StatusInternalServerError, "SERVER_ERROR", msg)
}

// SendUnauthorized 返回401错误, code=UNAUTHORIZED
func SendUnauthorized(w http.ResponseWriter, msg interface{}) {
	SendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

// SendForbidden 返回403错误, code=FORBIDDEN
func SendForbidden(w http.ResponseWriter, msg interface{}) {
	SendErrorResponse(w, http.StatusForbidden, "FORBIDDEN", msg)
}

// SendBusinessErrorAndCode 返回业务错误 httpCode=200, code=errCode参数
func SendBusinessErrorAndCode(w http.ResponseWriter, errCode string, msg interface{}) {
	SendErrorResponse(w, http.StatusOK, errCode, msg)
}

// SendSuccessResponse 返回成功结果
func SendSuccessResponse(w http.ResponseWriter, statusCode int, bizCode string, msg interface{}) {
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(BuildHttpResponse(bizCode, "T", msg))
}

// SendErrorResponse 返回失败结果
func SendErrorResponse(w http.ResponseWriter, statusCode int, errCode string, msg interface{}) {
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(BuildHttpResponse(errCode, "F", msg))
}

// BuildHttpResponse 构建返回json
func BuildHttpResponse(code string, flag string, str interface{}) []byte {
	bt, err := json.Marshal(HTTPResponse{Code: code, Flag: flag, Data: str})
	if nil != err {
		return []byte(err.Error())
	}
	return bt
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
