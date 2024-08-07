package gohttpx

import (
	"errors"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
)

func SetBaseUrl(base string) {
	baseURL = base
}

func GetBaseUrl() string {
	return baseURL
}

func AddGlobalHeader(k, v string) {
	globalHeaders[k] = v
}

func GetGlobalHeaders() map[string]string {
	return globalHeaders
}

// Headers["Accept"] = "application/json" for default
// 真正的请求逻辑
func doHttpRequest[T any](req *Request, t *T) *HttpError {
	if req.Url == "" || !strings.HasPrefix(req.Url, "http") {
		req.Url = GetBaseUrl() + req.Url
	}

	if req.Url == "" || !strings.HasPrefix(req.Url, "http") {
		return &HttpError{
			HttpStatusCode: HttpParamsError,
			Msg:            "url is invalid",
		}
	}

	if req.Timeout <= 0 {
		req.Timeout = time.Second * 10
	}

	var (
		restyClient = resty.New().
			SetTimeout(req.Timeout).
			SetRetryCount(req.RetryCount).
			SetRetryWaitTime(req.RetryWaitTime)
	)

	for k, v := range GetGlobalHeaders() {
		req.Headers[k] = v
	}

	if req.Headers == nil {
		req.Headers = make(map[string]string)
		req.Headers["Accept"] = CONTENT_TYPE_JSON
	} else {
		req.Headers["Accept"] = CONTENT_TYPE_JSON
	}

	var resp *resty.Response
	var err error
	request := restyClient.R()
	if req.Method == POST {
		request.
			SetBody(req.Body).
			SetQueryParams(req.QueryParams).
			SetResult(t).
			SetHeaders(req.Headers)

		if len(req.FormData) > 0 {
			request.SetFormData(req.FormData)
		}

		if len(req.Files) > 0 {
			request.SetFiles(req.Files)
		}
		resp, err = request.Post(req.Url)
	} else if req.Method == GET {
		resp, err = request.
			SetResult(t).
			SetQueryParams(req.QueryParams).
			SetHeaders(req.Headers).
			Get(req.Url)
	} else if req.Method == PUT {
		resp, err = request.
			SetResult(t).
			SetBody(req.Body).
			SetQueryParams(req.QueryParams).
			SetHeaders(req.Headers).
			Put(req.Url)
	} else {
		resp, err = request.
			SetResult(t).
			SetBody(req.Body).
			SetQueryParams(req.QueryParams).
			SetHeaders(req.Headers).
			Delete(req.Url)
	}

	if err != nil {
		return &HttpError{
			Error:          err,
			HttpStatusCode: HttpUnknownError,
			Msg:            "request timeout or unknown error",
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return &HttpError{
			Error:          err,
			HttpStatusCode: resp.StatusCode(),
			Msg:            "request timeout or unknown error",
		}
	}
	req.TraceInfo = resp.Request.TraceInfo() //调试信息
	respData, ok := resp.Result().(*T)
	if !ok {
		return &HttpError{
			Error:          err,
			HttpStatusCode: resp.StatusCode(),
			Msg:            "Response T is invalid",
		}
	}

	if respData == nil {
		return &HttpError{
			Error:          err,
			HttpStatusCode: resp.StatusCode(),
			Msg:            "Response data is empty",
		}
	}

	req.ResponseProto = resp.Proto()
	req.ResponseTime = resp.Time()
	req.Response = resp
	return nil
}

func HttpPostJson[T any](req *Request, t *T) *HttpError {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	req.Headers["Content-Type"] = CONTENT_TYPE_JSON
	req.Method = POST
	return HttpRequest[T](req, t)
}

func HttpPost[T any](req *Request, t *T) *HttpError {
	req.Method = POST
	return HttpRequest[T](req, t)
}

func HttpGet[T any](req *Request, t *T) *HttpError {
	req.Method = GET
	return HttpRequest[T](req, t)
}

// 带多个Urls重试逻辑
func HttpRequest[T any](req *Request, t *T) *HttpError {
	if req.IsConcurrency { //是并发开启
		return HttpConcurrencyRequest[T](req, t)
	}
	err := doHttpRequest[T](req, t)
	if err == nil {
		return nil
	} else {
		if len(req.Urls) == 0 { //没有重试urls
			return err
		}

		errs := &HttpError{
			HttpStatusCode: HttpRetryError,
			Error:          errors.New("HttpRetryError error"),
		}
		errs.Errors = append(errs.Errors, err)
		for _, url := range req.Urls {
			req.Url = url
			err = doHttpRequest[T](req, t)
			if err == nil { //请求成功了直接返回
				return nil
			} else {
				errs.Errors = append(errs.Errors, err) //请求失败继续,错误叠加记录
			}
		}
		return errs // 所有连接重试失败
	}
}

// 带多个Urls重试逻辑,并发请求,速度快先到达后 直接返回，其他请求取消
func HttpConcurrencyRequest[T any](req *Request, t *T) *HttpError {
	var err *HttpError
	if !req.IsAll {
		err = doHttpRequest[T](req, t)
		if err == nil {
			return nil
		}
	} else {
		if req.Url != "" { //把当前加进来起并发
			req.Urls = append(req.Urls, req.Url)
		}
	}

	if len(req.Urls) == 0 { //没有urls
		return err
	}

	errs := &HttpError{
		HttpStatusCode: HttpRetryError,
		Error:          errors.New("HttpRetryError error"),
	}
	if err != nil {
		errs.Errors = append(errs.Errors, err)
	}
	isSuccess := false
	fns := []func(){}
	for _, url := range req.Urls {
		reqNew := *req
		reqNew.Url = url
		fns = append(fns, func() {
			if isSuccess {
				return
			}
			err = doHttpRequest[T](&reqNew, t)
			if err != nil {
				errs.Errors = append(errs.Errors, err) //请求失败继续,错误叠加记录
			} else {
				isSuccess = true
				//请求成功了应该直接返回，剩下的请求结果忽略
			}
		})
	}

	goutils.AsyncFuncGroupOneSuccess(fns...)
	if isSuccess {
		return nil
	}
	return errs
}
