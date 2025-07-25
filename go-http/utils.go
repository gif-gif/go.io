package gohttp

import (
	"context"
	"github.com/gif-gif/go.io/go-utils/gocrypto"
	"github.com/gif-gif/go.io/go-utils/gozip"
	goserver "github.com/gif-gif/go.io/goio/server"
)

// POST 压缩请求
func CompressRequest(url string, body []byte, compressMethod string, compressType string, headers map[string]string) ([]byte, error) {
	_, data, err := gozip.Compress(body, compressMethod, compressType)
	if err != nil {
		return nil, err
	}

	payload := &Request{
		BinaryResponse: true,
		Url:            url,
		Method:         POST,
		Headers:        headers,
		Body:           data,
	}
	gh := GoHttp[goserver.Response]{
		Request: payload,
	}
	_, err = gh.HttpPost(context.Background())
	if err != nil {
		return nil, err
	}

	compressMethod = payload.Response.Header().Get(gozip.XNlContentEncoding)
	_, resData, err := gozip.Compress(payload.Response.Body(), compressMethod, gozip.UnGoZipType)
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// EncryptRequest 加密请求 aes cbc
func EncryptRequest(url string, body []byte, reqAesKey []byte, resAesKey []byte, compressMethod string, headers map[string]string) ([]byte, error) {
	data, err := gocrypto.GoDataEncrypt(body, reqAesKey, compressMethod)
	if err != nil {
		return nil, err
	}
	payload := &Request{
		BinaryResponse: true,
		Url:            url,
		Method:         POST,
		Headers:        headers,
		Body:           data,
	}
	gh := GoHttp[goserver.Response]{
		Request: payload,
	}
	_, err = gh.HttpPost(context.Background())
	if err != nil {
		return nil, err
	}

	compressMethod = payload.Response.Header().Get("X-NL-Content-Encoding")
	resData, err := gocrypto.GoDataDecrypt(payload.Response.Body(), resAesKey, compressMethod)
	if err != nil {
		return nil, err
	}
	return resData, nil
}

// aes ctr 加密请求
func EncryptCTRRequest(url string, body []byte, aesKey []byte, aesIv []byte, compressMethod string, headers map[string]string) ([]byte, error) {
	data, err := gocrypto.GoDataAesCTRTransformEncode(body, aesKey, aesIv, compressMethod)
	if err != nil {
		return nil, err
	}
	payload := &Request{
		BinaryResponse: true,
		Url:            url,
		Method:         POST,
		Headers:        headers,
		Body:           data,
	}
	gh := GoHttp[goserver.Response]{
		Request: payload,
	}
	_, err = gh.HttpPost(context.Background())
	if err != nil {
		return nil, err
	}

	compressMethod = payload.Response.Header().Get("X-NL-Content-Encoding")
	resData, err := gocrypto.GoDataAesCTRTransformDecode(payload.Response.Body(), aesKey, aesIv, compressMethod)
	if err != nil {
		return nil, err
	}
	return resData, nil
}
