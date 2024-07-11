package gohttpx

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

// Headers["Accept"] = "application/json" for default
func doHttpRequest[T any](req *Request, t *T) *HttpError {
	if req.Url == "" || !strings.HasPrefix(req.Url, "http") {
		return &HttpError{
			HttpStatusCode: HttpParamsError,
			Msg:            "url is invalid",
		}
	}
	var (
		restyClient = resty.New().
			SetTimeout(req.Timeout).
			SetRetryCount(req.RetryCount).
			SetRetryWaitTime(req.RetryWaitTime)
	)

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
		resp, err = request.
			SetBody(req.Body).
			SetQueryParams(req.QueryParams).
			SetResult(t).
			SetHeaders(req.Headers).
			Post(req.Url)
	} else {
		resp, err = request.
			SetResult(t).
			SetQueryParams(req.QueryParams).
			SetHeaders(req.Headers).
			Get(req.Url)
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
	err := doHttpRequest[T](req, t)
	if err == nil {
		return nil
	} else {
		if len(req.Urls) == 0 { //没有重试urls
			return err
		}

		errs := &HttpError{
			HttpStatusCode: HttpRetryError,
		}
		for _, url := range req.Urls {
			req.Url = url
			err = doHttpRequest[T](req, t)
			if err == nil { //请求成功了直接返回
				return err
			} else {
				errs.Errors = append(errs.Errors, err) //请求失败继续,错误叠加记录
			}
		}
		return errs
	}
}
