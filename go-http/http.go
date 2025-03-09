package gohttp

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

type GoHttp[T any] struct {
	BaseUrl string
	Headers map[string]string
	Request *Request
}

func (g *GoHttp[T]) SetBaseUrl(base string) {
	g.BaseUrl = base
}

func (g *GoHttp[T]) GetBaseUrl() string {
	return g.BaseUrl
}

func (g *GoHttp[T]) AddGlobalHeader(k, v string) {
	if g.Headers == nil {
		g.Headers = make(map[string]string)
	}
	g.Headers[k] = v
}

func (g *GoHttp[T]) RemoveGlobalHeader(k string) {
	delete(g.Headers, k)
}

func (g *GoHttp[T]) GetGlobalHeaders() map[string]string {
	return g.Headers
}

// Headers["Accept"] = "application/json" for default
// 真正的请求逻辑
func (g *GoHttp[T]) doHttpRequest(context context.Context, req *Request) (*T, error) {
	if req.Url == "" || !strings.HasPrefix(req.Url, "http") {
		req.Url = g.GetBaseUrl() + req.Url
	}

	if req.Url == "" || !strings.HasPrefix(req.Url, "http") {
		return nil, errors.New("[" + gconv.String(HttpParamsError) + "]" + "url is invalid")
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

	//if 2 {
	//	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	//}
	//
	//if 1 {
	//	restyClient.SetCertificates(certFile, keyFile, password)
	//}

	if req.proxyURL != "" {
		restyClient.SetProxy(req.proxyURL)
	}

	if req.Headers == nil {
		req.Headers = make(map[string]string)
		req.Headers["Accept"] = CONTENT_TYPE_JSON
	} else {
		req.Headers["Accept"] = CONTENT_TYPE_JSON
	}

	for k, v := range g.GetGlobalHeaders() {
		req.Headers[k] = v
	}

	var t T
	var resp *resty.Response
	var err error
	request := restyClient.R()
	request.SetContext(context)
	if !req.BinaryResponse {
		request.SetResult(t)
	}

	request.SetHeaders(req.Headers)
	if req.QueryParams != nil {
		request.SetQueryParams(req.QueryParams)
	}
	if req.ParamsValues != nil {
		request.SetQueryParamsFromValues(req.ParamsValues)
	}

	if req.Body != nil {
		request.SetBody(req.Body)
	}

	if req.Method == POST {
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
		resp, err = request.Get(req.Url)
	} else if req.Method == PUT {
		resp, err = request.Put(req.Url)
	} else {
		resp, err = request.Delete(req.Url)
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("[" + gconv.String(resp.StatusCode()) + "]" + "request timeout or unknown error->" + string(resp.Body()))
	}
	req.TraceInfo = resp.Request.TraceInfo() //调试信息
	if !req.BinaryResponse {
		respData, ok := resp.Result().(*T)
		if !ok {
			return nil, errors.New("[" + gconv.String(resp.StatusCode()) + "]" + "Response T is invalid")
		}

		if respData == nil {
			return nil, errors.New("[" + gconv.String(resp.StatusCode()) + "]" + "Response data is empty")
		}
	}

	req.ResponseProto = resp.Proto()
	req.ResponseTime = resp.Time()
	req.Response = resp
	return nil, nil
}

func (g *GoHttp[T]) HttpPostJson(context context.Context) (*T, error) {
	if g.Request.Headers == nil {
		g.Request.Headers = make(map[string]string)
	}

	g.Request.Headers["Content-Type"] = CONTENT_TYPE_JSON
	g.Request.Method = POST
	return g.HttpRequest(context)
}

func (g *GoHttp[T]) HttpPost(context context.Context) (*T, error) {
	g.Request.Method = POST
	return g.HttpRequest(context)
}

func (g *GoHttp[T]) HttpGet(context context.Context) (*T, error) {
	g.Request.Method = GET
	return g.HttpRequest(context)
}

// 带多个Urls重试逻辑
func (g *GoHttp[T]) HttpRequest(context context.Context) (*T, error) {
	t, err := g.doHttpRequest(context, g.Request)
	if err == nil {
		return t, nil
	} else {
		if len(g.Request.Urls) == 0 { //没有重试urls
			return t, err
		}

		errs := errors.New("HttpRetryError error")
		errs = errors.Join(errs, err)
		for _, url := range g.Request.Urls {
			g.Request.Url = url
			t, err = g.doHttpRequest(context, g.Request)
			if err == nil { //请求成功了直接返回
				return t, nil
			} else {
				errs = errors.Join(errs, err)
			}
		}
		return nil, errs // 所有连接重试失败
	}
}

// 带多个Urls重试逻辑,并发请求,速度快先到达后 直接返回，其他请求取消
func (g *GoHttp[T]) HttpConcurrencyRequest() (*T, error) {
	if g.Request.Url != "" { //把当前加进来起并发
		g.Request.Urls = append(g.Request.Urls, g.Request.Url)
	}

	if len(g.Request.Urls) == 0 { //没有urls
		return nil, errors.New("urls is empty")
	}

	var rst *T
	errs := errors.New("Concurrency error")
	var one sync.Once
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fns := []func(){}
	for _, url := range g.Request.Urls {
		reqNew := *g.Request
		reqNew.Url = url
		fns = append(fns, func() {
			if goutils.IsContextDone(ctx) {
				return
			}

			t, err := g.doHttpRequest(ctx, &reqNew)
			if err != nil {
				errs = errors.Join(errs, err)
			} else { //请求成功了应该直接返回，剩下的请求结果忽略
				one.Do(func() {
					rst = t
				})
				cancel() //有一个成功的取消所有请求
			}
		})
	}

	goutils.AsyncFuncGroup(fns...)
	if goutils.IsContextDone(ctx) {
		return rst, nil
	}

	return nil, errs
}
