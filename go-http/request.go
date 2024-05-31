package gohttp

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type Request struct {
	Headers map[string]string
	Tls     *Tls
	timeout time.Duration
	debug   bool
}

func (r *Request) Debug() *Request {
	r.debug = true
	return r
}

func (r *Request) SetHeader(name, value string) *Request {
	r.Headers[name] = value
	return r
}

func (r *Request) SetContentType(contentType string) *Request {
	r.SetHeader("Content-Type", contentType)
	return r
}

func (r *Request) JsonContentType() *Request {
	r.SetHeader("Content-Type", CONTENT_TYPE_JSON)
	return r
}

func (r *Request) SetTimeout(d time.Duration) *Request {
	r.timeout = d
	return r
}

func (r *Request) getClient() *http.Client {
	if r.timeout == 0 {
		r.timeout = 8 * time.Second
	}
	client := &http.Client{
		Timeout: r.timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if r.Tls != nil {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(r.Tls.CaCrt())
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{r.Tls.ClientCrt()},
			},
		}
	}
	return client
}

func (r *Request) Do(method, url string, reader io.Reader) (rst []byte, err error) {
	var (
		req *http.Request
		rsp *http.Response
	)

	req, err = http.NewRequest(method, url, reader)
	if err != nil {
		return
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	rsp, err = r.getClient().Do(req)
	if err != nil {
		return
	}

	defer rsp.Body.Close()

	rst, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	return
}

func (r *Request) handle(method, url string, data []byte) (rsp []byte, err error) {
	rsp, err = r.Do(method, url, bytes.NewReader(data))
	if r.debug {
		l := golog.WithTag(TAG).
			WithField("method", method).
			WithField("url", url).
			WithField("header", r.Headers).
			WithField("request-data", string(data)).
			WithField("response", string(rsp))
		if err != nil {
			l.Error()
		} else {
			l.Debug()
		}
	}
	return
}

func (r *Request) Get(url string) ([]byte, error) {
	return r.handle("GET", url, nil)
}

func (r *Request) GetWithQuery(url string, data []byte) ([]byte, error) {
	return r.handle("GET", url, data)
}

func (r *Request) Post(url string, data []byte) ([]byte, error) {
	return r.handle("POST", url, data)
}

func (r *Request) PostJson(url string, data []byte) ([]byte, error) {
	return r.JsonContentType().handle("POST", url, data)
}

func (r *Request) Put(url string, data []byte) ([]byte, error) {
	return r.handle("PUT", url, data)
}

func (r *Request) Upload(url, fileField, fileName string, fh io.Reader, data map[string]string) ([]byte, error) {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	goutils.AsyncFunc(func() {
		for k, v := range data {
			w.WriteField(k, v)
		}

		part, _ := w.CreateFormFile(fileField, fileName)

		io.CopyBuffer(part, fh, nil)

		w.Close()
		pw.Close()
	})

	r.SetHeader("Content-Type", w.FormDataContentType())

	return r.Do("POST", url, pr)
}
