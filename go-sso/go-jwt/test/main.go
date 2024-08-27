package main

import (
	golog "github.com/gif-gif/go.io/go-log"
	gojwt "github.com/gif-gif/go.io/go-sso/go-jwt"
)

func main() {
	c := gojwt.Config{
		AccessSecret:       "123456",
		AccessExpire:       100000,
		RefreshTokenExpire: 100000,
	}
	params := map[string]any{
		"aaa": 1,
	}

	gojwt := &gojwt.GoJwt{
		Config: c,
	}

	a, r, e, err := gojwt.GeneratedTokens(params)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	golog.Info(a, r, e)

	m, err := gojwt.ParseToken(a)
	if err != nil {
		golog.Error(err.Error())
		return
	}
	golog.Info(m)

	if gojwt.IsValidToken(a) {
		golog.WithTag("IsValidToken").Info("OK")
	}

	a, b, cc, err := gojwt.RefreshAccessToken(r)
	golog.WithTag("RefreshToken").Info(a, b, cc, err)
}
