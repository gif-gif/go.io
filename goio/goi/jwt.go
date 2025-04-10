package goi

import (
	gojwt "github.com/gif-gif/go.io/go-sso/go-jwt"
)

func GoJwt(names ...string) *gojwt.GoJwt {
	return gojwt.GetClient(names...)
}
