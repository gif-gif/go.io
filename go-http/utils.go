package gohttp

import (
	"context"
	"github.com/gif-gif/go.io/go-utils/gozip"
	goserver "github.com/gif-gif/go.io/goio/server"
)

// EncryptRequest 加密请求 aes cbc
func EncryptRequest(url string, body []byte, reqAesKey []byte, resAesKey []byte, compressMethod string, headers map[string]string) ([]byte, error) {
	data, err := gozip.GoDataEncrypt(body, reqAesKey, compressMethod)
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
	resData, err := gozip.GoDataDecrypt(payload.Response.Body(), resAesKey, compressMethod)
	if err != nil {
		return nil, err
	}
	return resData, nil
}

// aes ctr 加密请求
func EncryptCTRRequest(url string, body []byte, aesKey []byte, aesIv []byte, compressMethod string, headers map[string]string) ([]byte, error) {
	data, err := gozip.GoDataAesCTRTransformEncode(body, aesKey, aesIv, compressMethod)
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
	resData, err := gozip.GoDataAesCTRTransformDecode(payload.Response.Body(), aesKey, aesIv, compressMethod)
	if err != nil {
		return nil, err
	}
	return resData, nil
}
