package gooauth

import (
	"context"
	"golang.org/x/oauth2"
)

// OAuth2.0 授权接口
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
