package gohttpx

import (
	"bytes"
	"context"
	"errors"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/util/gconv"
	"net/http"
	"strings"
	"sync"
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

func RemoveGlobalHeader(k string) {
	delete(globalHeaders, k)
}

func GetGlobalHeaders() map[string]string {
	return globalHeaders
}

// Headers["Accept"] = "application/json" for default
// 真正的请求逻辑
func doHttpRequest[T any](context context.Context, req *Request, t *T) error {
	if req.Url == "" || !strings.HasPrefix(req.Url, "http") {
		req.Url = GetBaseUrl() + req.Url
	}

	if req.Url == "" || !strings.HasPrefix(req.Url, "http") {
		return errors.New("[" + gconv.String(HttpParamsError) + "]" + "url is invalid")
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

	if req.proxyURL != "" {
		restyClient.SetProxy(req.proxyURL)
	}

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
	request.SetContext(context)
	request.SetResult(t)
	if req.Method == POST {
		request.
			SetHeaders(req.Headers)

		if req.QueryParams != nil {
			request.SetQueryParams(req.QueryParams)
		}

		if req.Body != nil {
			request.SetBody(req.Body)
		}

		if len(req.FormData) > 0 {
			request.SetFormData(req.FormData)
		}

		if len(req.MultipartFormData) > 0 {
			request.SetMultipartFormData(req.MultipartFormData)
		}

		if len(req.FileBytes) > 0 {
			request.SetFileReader("file", req.FileName, bytes.NewReader(req.FileBytes))
		}

		if len(req.Files) > 0 {
			request.SetFiles(req.Files)
		}

		resp, err = request.Post(req.Url)
	} else if req.Method == GET {
		resp, err = request.
			SetQueryParams(req.QueryParams).
			SetHeaders(req.Headers).
			Get(req.Url)
	} else if req.Method == PUT {
		resp, err = request.
			SetBody(req.Body).
			SetQueryParams(req.QueryParams).
			SetHeaders(req.Headers).
			Put(req.Url)
	} else {
		resp, err = request.
			SetBody(req.Body).
			SetQueryParams(req.QueryParams).
			SetHeaders(req.Headers).
			Delete(req.Url)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("[" + gconv.String(resp.StatusCode()) + "]" + "request timeout or unknown error")
	}
	req.TraceInfo = resp.Request.TraceInfo() //调试信息
	respData, ok := resp.Result().(*T)
	if !ok {
		return errors.New("[" + gconv.String(resp.StatusCode()) + "]" + "Response T is invalid")
	}

	if respData == nil {
		return errors.New("[" + gconv.String(resp.StatusCode()) + "]" + "Response data is empty")
	}

	t = respData

	req.ResponseProto = resp.Proto()
	req.ResponseTime = resp.Time()
	req.Response = resp
	return nil
}

func HttpPostJson[T any](context context.Context, req *Request, t *T) error {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	req.Headers["Content-Type"] = CONTENT_TYPE_JSON
	req.Method = POST
	return HttpRequest[T](context, req, t)
}

func HttpPost[T any](context context.Context, req *Request, t *T) error {
	req.Method = POST
	return HttpRequest[T](context, req, t)
}

func HttpGet[T any](context context.Context, req *Request, t *T) error {
	req.Method = GET
	return HttpRequest[T](context, req, t)
}

// 带多个Urls重试逻辑
func HttpRequest[T any](context context.Context, req *Request, t *T) error {
	err := doHttpRequest[T](context, req, t)
	if err == nil {
		return nil
	} else {
		if len(req.Urls) == 0 { //没有重试urls
			return err
		}

		errs := errors.New("HttpRetryError error")
		errs = errors.Join(errs, err)
		for _, url := range req.Urls {
			req.Url = url
			err = doHttpRequest[T](context, req, t)
			if err == nil { //请求成功了直接返回
				return nil
			} else {
				errs = errors.Join(errs, err)
			}
		}
		return errs // 所有连接重试失败
	}
}

// 带多个Urls重试逻辑,并发请求,速度快先到达后 直接返回，其他请求取消
func HttpConcurrencyRequest[T any](req *Request, t *T) error {
	var err error
	if req.Url != "" { //把当前加进来起并发
		req.Urls = append(req.Urls, req.Url)
	}

	if len(req.Urls) == 0 { //没有urls
		return err
	}

	errs := errors.New("Concurrency error")

	var one sync.Once
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fns := []func(){}
	for _, url := range req.Urls {
		reqNew := *req
		reqNew.Url = url
		fns = append(fns, func() {
			if goutils.IsContextDone(ctx) {
				return
			}

			var tmp T
			err = doHttpRequest[T](ctx, &reqNew, &tmp)
			if err != nil {
				errs = errors.Join(errs, err)
			} else { //请求成功了应该直接返回，剩下的请求结果忽略
				one.Do(func() {
					t = &tmp
				})
				cancel() //有一个成功的取消所有请求
			}
		})
	}

	goutils.AsyncFuncGroup(fns...)
	if goutils.IsContextDone(ctx) {
		return nil
	}

	return errs
}
