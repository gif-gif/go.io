package gooauth

import (
	"context"
	"golang.org/x/oauth2"
)

// GoogleUserInfo 存储从 Google API 获取的用户信息
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// OAuth2.0 授权接口, 支持多种授权方式
type IGoOAuth interface {
	// 自定义参数
	//
	//	param1 := oauth2.SetAuthURLParam("param1", "value1")
	//	param2 := oauth2.SetAuthURLParam("param2", "value2")
	AuthUrl(opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, authorizationCode string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	RefreshToken(ctx context.Context, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	GetToken() *oauth2.Token
}
