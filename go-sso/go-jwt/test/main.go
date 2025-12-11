package main

import (
	golog "github.com/gif-gif/go.io/go-log"
	gojwt "github.com/gif-gif/go.io/go-sso/go-jwt"
)

func main() {
	c := gojwt.Config{
		AccessSecret:       "111",
		AccessExpire:       86400 * 360 * 10,
		RefreshTokenExpire: 86400 * 360 * 10,
	}

	params := map[string]any{
		"aaa": 1,
	}

	gojwt.Init(c)

	a, r, e, err := gojwt.Default().GeneratedTokens(params)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	golog.Info(a, r, e)

	m, ok, err := gojwt.Default().ParseTokenEx(a)
	if err != nil {
		golog.Error(err.Error())
		return
	}
	golog.Info(m, ok)

	if gojwt.Default().IsValidToken(a) {
		golog.WithTag("IsValidToken").Info("OK")
	}

	a, b, cc, err := gojwt.Default().RefreshAccessToken(r)
	golog.WithTag("RefreshToken").Info(a, b, cc, err)
}
