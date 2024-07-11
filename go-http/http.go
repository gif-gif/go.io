package gohttp

import (
	"io"
	"net/http"
	"strings"
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
