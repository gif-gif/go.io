package gohttpx

import (
	goutils "github.com/gif-gif/go.io/go-utils"
	"io"
	"net/http"
	"strings"
)

func GetHeaderIpInfos(r *http.Request, field string) map[string]string {
	result := make(map[string]string)
	result["X-Forwarded-For"] = r.Header.Get("X-Forwarded-For")
	result["X-Forward-For"] = r.Header.Get("X-Forward-For")
	result["Proxy-Client-IP"] = r.Header.Get("Proxy-Client-IP")
	result["WL-Proxy-Client-IP"] = r.Header.Get("WL-Proxy-Client-IP")
	result["HTTP_CLIENT_IP"] = r.Header.Get("HTTP_CLIENT_IP")
	result["HTTP_X_FORWARDED_FOR"] = r.Header.Get("HTTP_X_FORWARDED_FOR")
	result["X-Real-IP"] = r.Header.Get("X-Real-IP")
	result["r-RemoteAddr"] = r.RemoteAddr
	return result
}

// GetClientIp returns the client ip of this request without port.
// Note that this ip address might be modified by client header.
func GetClientIp(r *http.Request) string {
	clientIp := ""
	realIps := r.Header.Get("X-Forwarded-For")
	if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
		ipArray := strings.Split(realIps, ",")
		clientIp = ipArray[0]
		if clientIp != "" {
			if goutils.IsIPv4(clientIp) {
				return clientIp
			}
			if len(ipArray) > 1 {
				clientIp = strings.TrimSpace(ipArray[1])
				if goutils.IsIPv4(clientIp) {
					return clientIp
				}
			}
		}
	}

	realIps = r.Header.Get("X-Forward-For")
	if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
		ipArray := strings.Split(realIps, ",")
		clientIp = ipArray[0]
		if clientIp != "" {
			return clientIp
		}
	}

	if strings.EqualFold("unknown", realIps) {
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
		if idx := strings.LastIndex(clientIp, ":"); idx > 0 {
			clientIp = clientIp[:idx]
		}
	}

	return clientIp
}

func New(opts ...Option) *Request {
	r := &Request{
		Headers: map[string]string{
			"Content-Type": CONTENT_TYPE_FORM,
		},
	}
	for _, opt := range opts {
		switch opt.Name {
		case "tls":
			v := opt.Value.(map[string]string)
			r.Tls = &Tls{
				CaCrtFile:     v["caCrtFile"],
				ClientCrtFile: v["clientCrtFile"],
				ClientKeyFile: v["clientKeyFile"],
			}
		case "content-type-xml", "content-type-json", "content-type-form":
			r.SetHeader("Content-Type", opt.Value.(string))
		case "header":
			v := opt.Value.(map[string]string)
			for field, value := range v {
				r.SetHeader(field, value)
			}
		}
	}
	return r
}

func Get(url string) ([]byte, error) {
	return New().Get(url)
}

func GetWithQuery(url string, data []byte) ([]byte, error) {
	return New().GetWithQuery(url, data)
}

func Post(url string, data []byte) ([]byte, error) {
	return New().Post(url, data)
}

func PostJson(url string, data []byte) ([]byte, error) {
	return New().JsonContentType().Post(url, data)
}

func Put(url string, data []byte) ([]byte, error) {
	return New().Put(url, data)
}

func Upload(url, fileField, fileName string, fh io.Reader, data map[string]string) (b []byte, err error) {
	return New().Upload(url, fileField, fileName, fh, data)
}

func SetHeader(name, value string) *Request {
	return New().SetHeader(name, value)
}

func Debug() *Request {
	return New().Debug()
}
