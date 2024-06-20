package gohttpex

import (
	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/gogf/gf/util/gconv"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetClientIp returns the client ip of this request without port.
// Note that this ip address might be modified by client header.
func GetClientIp(r *http.Request) string {
	clientIp := ""
	realIps := r.Header.Get("X-Forwarded-For")
	if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
		ipArray := strings.Split(realIps, ",")
		clientIp = ipArray[0]
		if clientIp != "" {
			//fmt.Printf("GetClientIp X-Forwarded-For:%s\n", clientIp)
			return clientIp
		}
	}

	if clientIp == "" {
		realIps := r.Header.Get("X-Forward-For")
		if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
			ipArray := strings.Split(realIps, ",")
			clientIp = ipArray[0]
			if clientIp != "" {
				return clientIp
			}
		}
	}

	if clientIp == "" || strings.EqualFold("unknown", realIps) {
		clientIp = r.Header.Get("Proxy-Client-IP")
	}
	if clientIp == "" || strings.EqualFold("unknown", realIps) {
		clientIp = r.Header.Get("WL-Proxy-Client-IP")
	}
	if clientIp == "" || strings.EqualFold("unknown", realIps) {
		clientIp = r.Header.Get("HTTP_CLIENT_IP")
	}
	if clientIp == "" || strings.EqualFold("unknown", realIps) {
		clientIp = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if clientIp == "" || strings.EqualFold("unknown", realIps) {
		clientIp = r.Header.Get("X-Real-IP")
	}
	if clientIp == "" || strings.EqualFold("unknown", realIps) {
		clientIp = r.RemoteAddr
	}

	return clientIp
}

func HttpGetValuesResultBody(url string, params url.Values, headers map[string]string, retryCount int) ([]byte, error) {
	var (
		restyClient = resty.New().
			SetTimeout(time.Second * 20).
			EnableTrace().
			SetRetryCount(retryCount).
			SetRetryWaitTime(2 * time.Second)
	)

	resp, err := restyClient.R().
		//EnableTrace().
		SetQueryParamsFromValues(params).
		SetHeaders(headers).
		Get(url)

	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("请求错误" + gconv.String(resp.StatusCode()))
	}

	//ti := res.Request.TraceInfo()
	return resp.Body(), nil
}

func HttpGetValues[T any](url string, params url.Values, headers map[string]string, t *T, retryCount int) (*T, error) {
	var (
		restyClient = resty.New().
			SetTimeout(time.Second * 20).
			EnableTrace().
			SetRetryCount(retryCount).
			SetRetryWaitTime(2 * time.Second)
	)

	if headers == nil {
		headers = make(map[string]string)
		headers["Accept"] = "application/json"
	} else {
		headers["Accept"] = "application/json"
	}

	resp, err := restyClient.R().
		//EnableTrace().
		SetQueryParamsFromValues(params).
		SetResult(t).
		SetHeaders(headers).
		Get(url)

	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		//respData, ok := resp.Result().(*T)
		e1 := json.Unmarshal(resp.Body(), &t)
		if e1 != nil {
			return nil, errors.New("http Statusis not [200]OK : " + gconv.String(resp.StatusCode()) + ",body:" + string(resp.Body()))
		}
		return t, nil
	}

	//ti := res.Request.TraceInfo()

	respData, ok := resp.Result().(*T)
	if !ok {
		return nil, errors.New("params must be not empty ok is false")
	}
	if respData == nil {
		return nil, errors.New("data must be not empty")
	}

	return respData, nil
}

func HttpGet[T any](url string, params string, headers map[string]string, t *T, retryCount int) (*T, error) {
	var (
		restyClient = resty.New().
			SetTimeout(time.Second * 20).
			EnableTrace().
			SetRetryCount(retryCount).
			SetRetryWaitTime(2 * time.Second)
	)

	//value, _ := query.Values(OpenAiBalanceRequest{
	//	StartDate: req.StartDate,
	//	EndDate:   req.EndDate,
	//})SetFormDataFromValues
	if headers == nil {
		headers = make(map[string]string)
		headers["Accept"] = "application/json"
	} else {
		headers["Accept"] = "application/json"
	}

	resp, err := restyClient.R().
		//EnableTrace().
		SetQueryString(params).
		SetResult(t).
		SetHeaders(headers).
		Get(url)

	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		//respData, ok := resp.Result().(*T)
		e1 := json.Unmarshal(resp.Body(), &t)
		if e1 != nil {
			return nil, errors.New("http Statusis not [200]OK : " + gconv.String(resp.StatusCode()) + ",body:" + string(resp.Body()))
		}
		return t, nil
	}

	//ti := res.Request.TraceInfo()

	respData, ok := resp.Result().(*T)
	if !ok {
		return nil, errors.New("params must be not empty ok is false")
	}
	if respData == nil {
		return nil, errors.New("data must be not empty")
	}

	return respData, nil
}

func HttpPost[T any](url string, params map[string]interface{}, headers map[string]string, t *T) (*T, error) {
	var (
		restyClient = resty.New().
			SetTimeout(time.Second * 10).
			EnableTrace().
			SetRetryCount(0).
			SetRetryWaitTime(200 * time.Millisecond)
	)

	if headers == nil {
		headers = make(map[string]string)
		headers["Accept"] = "application/json"
	} else {
		headers["Accept"] = "application/json"
	}

	//value, _ := query.Values(OpenAiBalanceRequest{
	//	StartDate: req.StartDate,
	//	EndDate:   req.EndDate,
	//})SetFormDataFromValues

	resp, err := restyClient.R().
		//EnableTrace().
		SetBody(params).
		SetResult(t).
		SetHeaders(headers).
		Post(url)

	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("params must be not empty : " + gconv.String(resp.StatusCode()))
	}

	//ti := res.Request.TraceInfo()

	respData, ok := resp.Result().(*T)
	if !ok {
		return nil, errors.New("params must be not empty ok is false")
	}
	if respData == nil {
		return nil, errors.New("data must be not empty")
	}

	return respData, nil
}

func HttpPostJson[T any](url string, params string, headers map[string]string, timeout int, t *T) (*T, error) {
	var (
		restyClient = resty.New().
			SetTimeout(time.Second * time.Duration(timeout)).
			EnableTrace().
			SetRetryCount(0).
			SetRetryWaitTime(200 * time.Millisecond)
	)

	if headers == nil {
		headers = make(map[string]string)
		headers["Accept"] = "application/json"
	} else {
		headers["Accept"] = "application/json"
	}

	//value, _ := query.Values(OpenAiBalanceRequest{
	//	StartDate: req.StartDate,
	//	EndDate:   req.EndDate,
	//})SetFormDataFromValues
	headers["Content-Type"] = "application/json"
	resp, err := restyClient.R().
		//EnableTrace().
		SetBody(params).
		SetResult(t).
		SetHeaders(headers).
		Post(url)

	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("code:" + gconv.String(resp.StatusCode()))
	}

	//ti := res.Request.TraceInfo()

	respData, ok := resp.Result().(*T)
	if !ok {
		return nil, errors.New("params must be not empty ok is false")
	}
	if respData == nil {
		return nil, errors.New("data must be not empty")
	}

	return respData, nil
}
