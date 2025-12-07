package serviceutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wup364/pakku/pkg/constants/httpheaders"
	"github.com/wup364/pakku/pkg/constants/mediatypes"
	"github.com/wup364/pakku/pkg/serviceutil/formutil"
	"github.com/wup364/pakku/pkg/utypes"
)

// ErrUnknownError 未知错误
var ErrUnknownError = utypes.NewCustomError(errors.New("unknown error"), "SERVER_ERROR")

// HandleRequest 根据请求类型将请求参数转换为对象并执行处理函数
// 执行处理函数: h = func(h T)(any, error) 或 func(h T)error
func HandleRequest[T any](r *http.Request, w http.ResponseWriter, h any) {
	// 解析参数
	var cmd T
	if err := ParseHTTPRequest(r, &cmd); nil != err {
		SendServerError(w, err.Error())
		return
	}

	// 执行
	if fn, ok := h.(func(h T) (any, error)); ok {
		HandleRequestWithData(w, func() (v any, err error) {
			return fn(cmd)
		})

	} else if fn, ok := h.(func(h T) error); ok {
		HandleRequestWithoutData(w, func() (err error) {
			return fn(cmd)
		})

	} else {
		SendServerError(w, ErrUnknownError.Error())
	}
}

// HandleRequestWithoutData 执行并自动响应(无响应数据的)
func HandleRequestWithoutData(w http.ResponseWriter, fun func() (err error)) {
	if err := fun(); nil != err {
		SendBusinessError(w, err)
	} else {
		SendSuccess(w, "")
	}
}

// HandleRequestWithData 执行并自动响应(有响应数据的)
func HandleRequestWithData(w http.ResponseWriter, fun func() (v any, err error)) {
	if res, err := fun(); nil != err {
		SendBusinessError(w, err)
	} else {
		SendSuccess(w, res)
	}
}

// ParseHTTPRequest 根据请求类型将请求参数转换为对象
func ParseHTTPRequest(r *http.Request, obj any) error {
	contentType := r.Header.Get(httpheaders.CONTENT_TYPE)

	switch {
	case strings.Contains(contentType, mediatypes.APPLICATION_JSON):
		return ParseJSONRequest(r, obj)
	case strings.Contains(contentType, mediatypes.APPLICATION_FORM_URLENCODED):
		return ParseFormRequest(r, obj)
	default:
		// 在 default 分支中，根据请求方法选择解析方式
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			return ParseFormRequest(r, obj)
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect:
			return ParseJSONRequest(r, obj)
		default:
			return fmt.Errorf("unsupported content type or method: %s, %s", contentType, r.Method)
		}
	}
}

// GetRequestRange 解析http分段头信息
func GetRequestRange(r *http.Request, maxSize int64) (start, end int64, hasRange bool) {
	var qRange string
	if qRange = r.Header.Get(httpheaders.RANGE); len(qRange) == 0 {
		qRange = r.FormValue(httpheaders.RANGE)
	}

	maxEnd := maxSize - 1
	if len(qRange) > 0 {
		hasRange = true
		temp := qRange[strings.Index(qRange, "=")+1:]
		if index := strings.Index(temp, "-"); index > -1 {
			var err error
			if start, err = strconv.ParseInt(temp[0:strings.Index(temp, "-")], 10, 64); nil != err || start < 0 {
				start = 0
			}
			if end, err = strconv.ParseInt(temp[strings.Index(temp, "-")+1:], 10, 64); nil != err || end == 0 {
				end = maxEnd
			}
		}
	} else {
		end = maxEnd
	}
	return start, end, hasRange
}

// ParseJSONRequest 解析 http 请求的 JSON 数据并绑定到对象上
func ParseJSONRequest(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

// ParseFormRequest 解析 http 请求的 form 参数并绑定到对象上
func ParseFormRequest(r *http.Request, obj any) error {
	if r.Form == nil && r.PostForm == nil {
		if err := r.ParseForm(); err != nil {
			return err
		}
	}
	return formutil.NewFormBinder().Bind(r, obj)
}
